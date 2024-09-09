package handler

import "github.com/bootstrap-library/stock-plus/telegram"

var Handlers = make(map[string]telegram.MessageHandler)

func Register(mode string, handler telegram.MessageHandler) {
	Handlers[mode] = handler
}

func GetHandler(mode string) telegram.MessageHandler {
	if handler, ok := Handlers[mode]; ok {
		return handler
	}

	return nil
}
