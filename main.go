package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	list := tview.NewList()
	list.SetBorder(true).SetTitle("Chat List <Ctrl-L>")

	go func() {
		chatList, err := db.GetChatList()
		if err != nil {
			log.Fatal(err)
		}
		for _, chat := range chatList {
			if chat.ChatTitle == "" {
				continue
			}
			list.AddItem(chat.ChatTitle, fmt.Sprintf("ChatId:%d", chat.ChatID), rune(0), nil)
		}
	}()
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			// 渲染history box
			go func() {
				_, secondary := list.GetItemText(list.GetCurrentItem())
				app.SetFocus(textArea)
				// db.GetChatHistoryByChatId()
				secs := strings.Split(secondary, ":")
				chatID, _ := strconv.ParseInt(secs[1], 10, 64)
				history, err := db.GetChatHistoryByChatId(chatID)
				if err != nil {
					log.Fatal(err)
				}
				for _, h := range history {
					if h.Role == "user" {
						historyChan <- "[red]Q:" + h.Prompt + "\n"
					} else if h.Role == "model" {
						historyChan <- "[green]A:" + h.Prompt + "\n"
					}
				}
			}()

			return nil
		}
		return event
	})
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatLog, 0, 8, false).
		AddItem(textArea, 0, 2, false).
		AddItem(helpInfo, 1, 0, false)

	appFlex := tview.NewFlex().AddItem(list, 0, 3, true).AddItem(flex, 0, 7, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlH:
			app.SetFocus(chatLog)
			return nil
		case tcell.KeyCtrlI:
			app.SetFocus(textArea)
			return nil
		case tcell.KeyCtrlL:
			app.SetFocus(list)
			return nil
		}
		return event
	})

	if err := app.SetRoot(appFlex, true).Run(); err != nil {
		panic(err)
	}
}
