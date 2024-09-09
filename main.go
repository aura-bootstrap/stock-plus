package main

import (
	"bytes"
	"fmt"

	"github.com/bootstrap-library/stock-plus/config"
	"github.com/bootstrap-library/stock-plus/env"
	"github.com/bootstrap-library/stock-plus/telegram"
)

func main() {
	start(messageHandler)
}

func start(handler telegram.MessageHandler) {
	telegramServer := &telegram.Server{
		Config: telegram.Config{
			APIToken:  config.TelegramAPITokenList[env.Int("INSTANCE")],
			FirstName: config.TelegramFirstName,
			Proxy:     config.TelegramProxy,
		},
	}
	telegramServer.Init()
	telegramServer.Serve(handler)
}

func messageHandler(text string, sender telegram.MessageSender) {
	sender(display())
}

var Items = []string{"延迟测试"}

func display() string {
	var b bytes.Buffer
	b.WriteString("当前支持的功能有：\n")
	for i, item := range Items {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}
	return b.String()
}
