package main

import (
	"log"
	"sync"
	"tg_bot_id/internal/config"
	"tg_bot_id/internal/server"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var waitGo sync.WaitGroup

func main() {

	waitGo.Add(2)

	go func() {
		defer waitGo.Done()
		server.StartWebServer()
	}()

	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("[ERROR] failed to create bot %v", err)
		return
	}

	botAPI.Debug = true
	log.Printf("Authorized on account %v\n\n", botAPI.Self.UserName)

	go func() {
		defer waitGo.Done()
		server.StartBot(botAPI)
	}()

	waitGo.Wait()
}
