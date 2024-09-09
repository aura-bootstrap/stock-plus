package chat

import (
	"github.com/bootstrap-library/stock-plus/handler"
	"github.com/bootstrap-library/stock-plus/telegram"
)

var bot *ChatBot

func init() {
	handler.Register("chat", func(text string, sender telegram.MessageSender) {
		if bot == nil {
			bot = NewChatBot(sender)
		}
		bot.HandleInput(text)
	})
}

// ChatBot 是状态机结构体
type ChatBot struct {
	currentState State
	sender       telegram.MessageSender
}

func NewChatBot(sender telegram.MessageSender) *ChatBot {
	return &ChatBot{
		currentState: &StartState{},
		sender:       sender,
	}
}

func (c *ChatBot) SetState(state State) {
	c.currentState = state
}

func (c *ChatBot) HandleInput(input string) {
	c.currentState = c.currentState.HandleInput(input)
}

func (c *ChatBot) Display() string {
	return c.currentState.Display()
}
