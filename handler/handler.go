package handler

import (
	"fmt"

	"github.com/bootstrap-library/stock-plus/telegram"
)

var currentHandler telegram.MessageHandler

func BootHandler(text string, sender telegram.MessageSender) {
	if text == "boot" {
		currentHandler = nil
		sender("I'm Bootstrap Bot, please input a mode")
		return
	}

	handler := GetHandler(text)
	if handler == nil {
		sender("Please input a valid mode")
		return
	}

	currentHandler = handler
	sender(fmt.Sprintf("Enter [%s] mode", text))
	return
}
