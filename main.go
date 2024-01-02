package main

import (
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var HOME_PATH = os.Getenv("HOME")
var chatID = 1

func main() {

	sendMsgChan := make(chan string)
	historyChan := make(chan string)

	geminiClient, err := newGeminiClient()
	if err != nil {
		log.Fatal(err)
	}

	go geminiClient.sendMessageToTui(sendMsgChan, historyChan, db)

	geminiTui := NewGeminiTui(historyChan, sendMsgChan)
	go func() {
		for {
			history := <-historyChan
			geminiTui.TuiApp.QueueUpdateDraw(func() {
				geminiTui.ChatHistoryUI.Write([]byte(history))
			})
		}
	}()

	geminiTui.Run()
}
