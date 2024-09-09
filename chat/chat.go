package chat

import (
	"github.com/bootstrap-library/stock-plus/handler"
	"github.com/bootstrap-library/stock-plus/telegram"
)

var bot *ChatBot

func init() {
	handler.Register("chat", func(text string, sender telegram.MessageSender) {
		if bot == nil && text == "/chat" {
			bot = NewChatBot(sender)
			return
		}
		bot.Handle(text)
	})
}

// State 接口定义了所有状态必须实现的方法
type State interface {
	String() string
	EnterState() string
	HandleInput(input string) (State, string)
	LeaveState() string
}

// ChatBot 是状态机结构体
type ChatBot struct {
	currentState State
	sender       telegram.MessageSender
}

func NewChatBot(sender telegram.MessageSender) *ChatBot {
	bot := &ChatBot{
		sender: sender,
	}
	bot.ChangeState(&MenuState{})
	return bot
}

func (c *ChatBot) Handle(input string) {
	newState, output := c.currentState.HandleInput(input)
	if output != "" {
		c.sender(output)
	}
	if newState != nil && newState.String() != c.currentState.String() {
		c.ChangeState(newState)
	}
}

func (c *ChatBot) ChangeState(state State) {
	if c.currentState != nil {
		output := c.currentState.LeaveState()
		if output != "" {
			c.sender(output)
		}
	}
	c.currentState = state
	if c.currentState != nil {
		output := c.currentState.EnterState()
		if output != "" {
			c.sender(output)
		}
	}
}
