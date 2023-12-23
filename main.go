package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\x1b[31mQ: ")
		text, _ := reader.ReadString('\n')
		fmt.Print("\x1b[0m")
		text = strings.Replace(text, "\n", "", -1)
		geminiClient.sendMessageStreamAndPrint(text, db)
	}
}
