package handler

import (
	"fmt"
	"strings"

	"github.com/bootstrap-library/stock-plus/telegram"
)

var currentHandler telegram.MessageHandler

func BootHandler(text string, sender telegram.MessageSender) {
	if strings.HasPrefix(text, "/") {
		handler := GetHandler(text[1:])
		if handler == nil {
			sender("Please input a valid mode.")
			return
		}

		currentHandler = handler
		sender(fmt.Sprintf("Enter %s mode.", text))
	}

	currentHandler(text, sender)
}
