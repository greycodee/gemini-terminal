package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var HOME_PATH = os.Getenv("HOME")
var chatID = 1

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

	geminiTui := NewGeminiTui(historyChan, sendMsgChan, db)
	go func() {
		for {
			history := <-historyChan
			geminiTui.TuiApp.QueueUpdateDraw(func() {
				geminiTui.ChatHistoryUI.Write([]byte(history))
			})
		}
	}()

	go func() {
		title := <-titleChan
		geminiTui.TuiApp.QueueUpdateDraw(func() {
			geminiTui.ChatHistoryUI.SetTitle(geminiTui.ChatHistoryUI.GetTitle() + " [" + title + "]")
		})
	}()

	go func() {
		for {
			flag := <-genFlagChan
			if flag {
				geminiTui.TuiApp.SetFocus(geminiTui.ChatHistoryUI)
			} else {
				geminiTui.TuiApp.SetFocus(geminiTui.ChatInputUI)
			}
		}
	}()

	go func() {
		chatList, err := db.GetChatList()
		if err != nil {
			log.Fatal(err)
		}
		for _, chat := range chatList {
			if chat.ChatTitle == "" {
				continue
			}
			geminiTui.ChatListUI.AddItem(chat.ChatTitle, fmt.Sprintf("ChatId:%d", chat.ChatID), rune(0), nil)
		}
	}()
	geminiTui.Run()
}
