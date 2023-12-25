package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	_ "github.com/mattn/go-sqlite3"
)

var HOME_PATH = os.Getenv("HOME")
var chatID = 1
var config GeminiChatConfig

func init() {
	// 获取参数
	if len(os.Args) > 1 {
		chatIDStr := os.Args[1]
		chatIDTmerl, err := strconv.Atoi(chatIDStr)
		if err != nil {
			log.Fatal(err)
		}
		chatID = chatIDTmerl
	}
	config = GetConfig()
}

func main() {
	sendMsgChan := make(chan string)
	historyChan := make(chan string)
	genFlagChan := make(chan bool)

	app := tview.NewApplication()

	ctx := context.Background()
	db := initDB()
	defer db.SqliteDB.Close()
	geminiClient, err := newGeminiClient(ctx, chatID, config)
	if err != nil {
		log.Fatal(err)
	}

	history, err := db.GetByChatID(chatID)
	if err != nil {
		log.Fatal(err)
	}
	geminiClient.startChat(history)
	defer geminiClient.client.Close()

	go geminiClient.sendMessageToTui(sendMsgChan, historyChan, genFlagChan, db)

	chatLog := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	chatLog.SetChangedFunc(func() {
		app.Draw()
		chatLog.ScrollToEnd()
	})
	go func() {
		for {
			history := <-historyChan
			app.QueueUpdateDraw(func() {
				chatLog.Write([]byte(history))
				// chatLog.SetText(chatLog.GetText(true) + history)
			})
		}
	}()

	// 创建一个输入框用于输入消息
	inputField := tview.NewInputField().
		SetLabel("Enter message: ").
		SetFieldWidth(0)

		// inputField.SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).SetEnabled(false)

	// 设置完成函数以处理消息提交
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			sendMsgChan <- inputField.GetText()
			inputField.SetText("")
			genFlagChan <- true
		}
	})

	go func() {
		for {
			flag := <-genFlagChan
			inputField.SetDisabled(flag)
		}
	}()
	// spinner := tview.NewSpinner('⣾', '⣽', '⣻', '⢿', '⡿', '⣟', '⣯', '⣷')
	// 创建一个Flex布局，并设置聊天记录框和输入框的比例为8:1
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatLog, 0, 8, false).
		AddItem(inputField, 1, 2, true)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
