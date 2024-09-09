package telegram

import (
	"log"
	"net/http"
	"net/url"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
)

type MessageSender func(string)
type MessageHandler func(string, MessageSender)

var bot *api.BotAPI

type Config struct {
	APIToken  string
	FirstName string
	Proxy     string
}

type Server struct {
	Config
}

func (s *Server) Init() {
	var err error
	bot, err = api.NewBotAPIWithClient(s.APIToken, api.APIEndpoint, s.newClient())
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func (s *Server) Serve(messageHandler MessageHandler) {
	u := api.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Chat.FirstName == s.FirstName {
			log.Printf("[Bot] <- %s", update.Message.Text)

			messageHandler(update.Message.Text, func(text string) {
				msg := api.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg)

				log.Printf("[Bot] -> %s", text)
			})
		}
	}
}

func (s *Server) newClient() *http.Client {
	c := &http.Client{
		Transport: &http.Transport{},
	}

	u, err := url.Parse(s.Proxy)
	if err != nil {
		log.Panic(err)
	}

	if u.Scheme == "socks5" {
		auth := proxy.Auth{}
		if u.User != nil {
			auth.User = u.User.Username()
			auth.Password, _ = u.User.Password()
		}
		dialer, err := proxy.SOCKS5("tcp", u.Host, &auth, proxy.Direct)
		if err != nil {
			log.Panic(err)
		}
		c.Transport.(*http.Transport).DialContext = dialer.(proxy.ContextDialer).DialContext
	} else {
		c.Transport.(*http.Transport).Proxy = http.ProxyURL(u)
	}

	return c
}
