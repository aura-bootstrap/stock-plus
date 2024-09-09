package main

import (
	"github.com/bootstrap-library/stock-plus/config"
	"github.com/bootstrap-library/stock-plus/env"
	"github.com/bootstrap-library/stock-plus/handler"
	"github.com/bootstrap-library/stock-plus/telegram"
)

func main() {
	telegramServer := &telegram.Server{
		Config: telegram.Config{
			APIToken:  config.TelegramAPITokenList[env.Int("INSTANCE")],
			FirstName: config.TelegramFirstName,
			Proxy:     config.TelegramProxy,
		},
	}
	telegramServer.Init()
	telegramServer.Serve(handler.BootHandler)
}
