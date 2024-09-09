package main

import (
	"github.com/bootstrap-library/stock-plus/config"
	"github.com/bootstrap-library/stock-plus/env"
	"github.com/bootstrap-library/stock-plus/handler"
	"github.com/bootstrap-library/stock-plus/telegram"

	_ "github.com/bootstrap-library/stock-plus/chat"
)

func main() {
	telegramServer := &telegram.Server{
		Config: telegram.Config{
			APIToken:  config.TelegramAPITokenList[env.Int("INSTANCE")-1],
			FirstName: config.TelegramFirstName,
			Proxy:     config.TelegramProxy,
		},
	}
	telegramServer.Init()
	telegramServer.Serve(handler.BootHandler)
}
