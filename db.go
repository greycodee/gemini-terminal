package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
)

const (
	DB_PATH = "/.local/share/gemini/"
	DB_NAME = "gemini.db"
)

type GeminiChatHistory struct {
	ID         int64  `db:"id"`
	ChatID     int64  `db:"chat_id"`
	Prompt     string `db:"prompt"`
	Role       string `db:"role"`
	CreateTime string `db:"create_time"`
}

type DB struct {
	SqliteDB *sql.DB
}

func initDB() *DB {
	FULL_DB_PATH := HOME_PATH + DB_PATH
	if _, err := os.Stat(FULL_DB_PATH); os.IsNotExist(err) {
		os.MkdirAll(FULL_DB_PATH, os.ModePerm)
	}

	sqliteDB, err := sql.Open("sqlite3", FULL_DB_PATH+DB_NAME)
	if err != nil {
		log.Fatal(err)
	}
	_, err = sqliteDB.Exec(`CREATE TABLE IF NOT EXISTS gemini_chat_history (
		id INTEGER PRIMARY KEY,
		chat_id INTEGER,
		prompt TEXT,
		role TEXT,
		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		SqliteDB: sqliteDB,
	}
}

func (db *DB) InsertHistory(chat GeminiChatHistory) error {
	_, err := db.SqliteDB.Exec(`INSERT INTO gemini_chat_history (chat_id, prompt, role) VALUES (?, ?, ?)`, chat.ChatID, chat.Prompt, chat.Role)
	return err
}

func (db *DB) InsertHistoryWithTX(tx *sql.Tx, chat GeminiChatHistory) error {
	_, err := tx.Exec(`INSERT INTO gemini_chat_history (chat_id, prompt, role) VALUES (?, ?, ?)`, chat.ChatID, chat.Prompt, chat.Role)
	return err
}

func (db *DB) GetByChatID(chatId int) ([]*genai.Content, error) {
	chatHistoryList := make([]*genai.Content, 0)
	rows, err := db.SqliteDB.Query(`SELECT prompt,role FROM gemini_chat_history WHERE chat_id = ?`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		var prompt string
		var role string
		err := rows.Scan(&prompt, &role)
		if err != nil {
			return nil, err
		}
		chatHistoryList = append(chatHistoryList, &genai.Content{
			Parts: parsePrompt(prompt),
			Role:  role,
		})
	}
	return chatHistoryList, nil
}

func parsePrompt(prompt string) []genai.Part {
	// 解析prompt数组
	var promptList []string
	err := json.Unmarshal([]byte(prompt), &promptList)
	if err != nil {
		log.Fatal(err)
	}
	promptPart := make([]genai.Part, 0)
	for _, v := range promptList {
		promptPart = append(promptPart, genai.Text(v))
	}
	return promptPart
}
