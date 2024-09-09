package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/bootstrap-library/stock-plus/crypto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
)

var (
	APIToken  = crypto.Decrypt("4rahgCTNrFjpMLibS3cu4w6PP//JsoWwdxDWPDzEUB/eMFw6pjtsb0AV0CfZtVYB")
	FirstName = crypto.Decrypt("YfS9C5G3V9loLwij1ZSBpwdTimK1Nz5chM1FjnVSCIA=")
)

func main1() {
	socks5Proxy := "socks5://127.0.0.1:1080"
	proxyURL, err := url.Parse(socks5Proxy)
	if err != nil {
		log.Fatalf("Failed to parse proxy URL: %v", err)
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		log.Fatalf("Failed to create proxy dialer: %v", err)
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	bot, err := tgbotapi.NewBotAPIWithClient(APIToken, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		log.Panic(err)
	}
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Chat.FirstName == FirstName {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
