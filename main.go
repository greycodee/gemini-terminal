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
var chatHistoryTitle = "Chat History <Ctrl-H>"

func main() {
	db := initDB()
	defer db.SqliteDB.Close()
	// 获取参数
	if len(os.Args) > 1 {
		chatIDStr := os.Args[1]
		chatIDTmerl, err := strconv.Atoi(chatIDStr)
		if err != nil {
			log.Fatal(err)
		}
		chatID = chatIDTmerl
	} else {
		id, err := db.GetLatestChatID()
		if err != nil {
			log.Fatal(err)
		}
		chatID = id + 1
	}
	config := GetConfig()

	sendMsgChan := make(chan string)
	historyChan := make(chan string)
	genFlagChan := make(chan bool)
	titleChan := make(chan string)

	app := tview.NewApplication()

	ctx := context.Background()

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

	go geminiClient.sendMessageToTui(sendMsgChan, historyChan, genFlagChan, titleChan, db)

	chatLog := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	chatLog.SetChangedFunc(func() {
		app.Draw()
		chatLog.ScrollToEnd()
	}).SetBorder(true).SetTitle(chatHistoryTitle)

	go func() {
		for {
			history := <-historyChan
			app.QueueUpdateDraw(func() {
				chatLog.Write([]byte(history))
			})
		}
	}()

	go func() {
		title := <-titleChan
		app.QueueUpdateDraw(func() {
			chatLog.SetTitle(chatHistoryTitle + " [" + title + "]")
		})
	}()

	textArea := tview.NewTextArea()
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlS {
			if textArea.GetText() == "" {
				return nil
			}
			sendMsgChan <- textArea.GetText()
			genFlagChan <- true
			textArea.SetText("", true)
			return nil
		}
		return event
	})
	textArea.SetBorder(true).SetTitle("Input Message <Ctrl-I>")

	go func() {
		for {
			flag := <-genFlagChan
			if flag {
				app.SetFocus(chatLog)
			} else {
				app.SetFocus(textArea)
			}
		}
	}()

	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-S to send message, press Ctrl-C to exit")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatLog, 0, 8, false).
		AddItem(textArea, 0, 2, true).
		AddItem(helpInfo, 1, 0, true)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlH:
			app.SetFocus(chatLog)
			return nil
		case tcell.KeyCtrlI:
			app.SetFocus(textArea)
			return nil
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
