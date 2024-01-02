package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GeminiTui struct {
	TuiApp          *tview.Application
	ChatListUI      *tview.List
	ChatHistoryUI   *tview.TextView
	ChatInputUI     *tview.TextArea
	TuiFlexUI       *tview.Flex
	ChatHistoryChan chan string
	SendMsgChan     chan string
	db              *DB
}

func NewGeminiTui(chatHistoryChan, sendMsgChan chan string) GeminiTui {
	// TODO init db
	// TODO init geminiClient
	db := initDB()
	geminiClient, err := newGeminiClient()
	if err != nil {
		log.Fatal(err)
	}

	history, err := db.GetByChatID(chatID)
	if err != nil {
		log.Fatal(err)
	}
	geminiClient.startChat(history)

	tuiApp := tview.NewApplication()
	geminiTui := GeminiTui{
		TuiApp:          tuiApp,
		ChatHistoryChan: chatHistoryChan,
		SendMsgChan:     sendMsgChan,
		db:              db,
	}
	geminiTui.TuiApp.SetInputCapture(geminiTui.moveWindowsFocus)
	geminiTui.drawTui()

	// TODO init data
	geminiTui.initData()
	return geminiTui
}

func (tui *GeminiTui) initData() {
	go func() {
		chatList, err := tui.db.GetChatList()
		if err != nil {
			log.Fatal(err)
		}
		for _, chat := range chatList {
			if chat.ChatTitle == "" {
				continue
			}
			tui.ChatListUI.AddItem(chat.ChatTitle, fmt.Sprintf("ChatId:%d", chat.ChatID), rune(0), nil)
		}
	}()
}

func (tui *GeminiTui) Run() {
	if err := tui.TuiApp.SetRoot(tui.TuiFlexUI, true).Run(); err != nil {
		panic(err)
	}
}

func (tui *GeminiTui) drawTui() {
	leftHelpInfo := tview.NewTextView().
		SetText(" Press Ctrl-N to create new Chat Session.")
	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-S to send message, press Ctrl-C to exit")
	tui.genChatHistoryUI()
	tui.genChatInputUI()
	tui.genChatListUI()

	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.ChatListUI, 0, 1, true).
		AddItem(leftHelpInfo, 1, 1, false)

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.ChatHistoryUI, 0, 8, false).
		AddItem(tui.ChatInputUI, 0, 2, false).
		AddItem(helpInfo, 1, 0, false)

	appFlex := tview.NewFlex().
		AddItem(leftFlex, 0, 5, true).
		AddItem(rightFlex, 0, 7, false)
	tui.TuiFlexUI = appFlex
}

func (tui *GeminiTui) genChatHistoryUI() {
	chatHistoryUI := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	chatHistoryUI.
		SetBorder(true).
		SetTitle("Chat History <Ctrl-H>")
	tui.ChatHistoryUI = chatHistoryUI
	tui.ChatHistoryUI.SetChangedFunc(func() {
		tui.TuiApp.Draw()
		tui.ChatHistoryUI.ScrollToEnd()
	})
}

func (tui *GeminiTui) genChatListUI() {
	chatListUI := tview.NewList()
	chatListUI.
		SetBorder(true).
		SetTitle("Chat List <Ctrl-L>")
	tui.ChatListUI = chatListUI
	tui.ChatListUI.SetInputCapture(tui.handlerChatListUIKeyEvent)
}

func (tui *GeminiTui) genChatInputUI() {
	inputArea := tview.NewTextArea()
	inputArea.
		SetBorder(true).
		SetTitle("Input Message <Ctrl-I>")
	tui.ChatInputUI = inputArea
	tui.ChatInputUI.SetInputCapture(tui.handlerChatInputUIKeyEvent)
}

func (tui *GeminiTui) moveWindowsFocus(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlH:
		tui.TuiApp.SetFocus(tui.ChatHistoryUI)
		return nil
	case tcell.KeyCtrlI:
		tui.TuiApp.SetFocus(tui.ChatInputUI)
		return nil
	case tcell.KeyCtrlL:
		tui.TuiApp.SetFocus(tui.ChatListUI)
		return nil
	default:
		return event
	}
}

func (tui *GeminiTui) handlerChatListUIKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlN:
		tui.ChatListUI.AddItem("New Chat Session", fmt.Sprintf("ChatId:%d", chatID), rune(0), nil)
		tui.TuiApp.SetFocus(tui.ChatInputUI)
		return nil
	case tcell.KeyEnter:
		go func() {
			_, secondary := tui.ChatListUI.GetItemText(tui.ChatListUI.GetCurrentItem())
			tui.TuiApp.SetFocus(tui.ChatInputUI)
			// db.GetChatHistoryByChatId()
			secs := strings.Split(secondary, ":")
			chatID, _ := strconv.ParseInt(secs[1], 10, 64)
			history, err := tui.db.GetChatHistoryByChatId(chatID)
			if err != nil {
				log.Fatal(err)
			}
			for _, h := range history {
				if h.Role == "user" {
					tui.ChatHistoryChan <- "[red]Q:" + h.Prompt + "\n"
				} else if h.Role == "model" {
					tui.ChatHistoryChan <- "[green]A:" + h.Prompt + "\n"
				}
			}
		}()
		return nil
	default:
		return event
	}
}

func (tui *GeminiTui) handlerChatInputUIKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlS:
		// TODO publish message
		if tui.ChatInputUI.GetText() == "" {
			return nil
		}
		tui.SendMsgChan <- tui.ChatInputUI.GetText()
		tui.ChatInputUI.SetText("", true)
		return nil
	default:
		return event
	}
}
