package main

import (
	"log"
	"tg_bot_id/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("[ERROR] failed to create bot %v", err)
		return
	}

	botAPI.Debug = true
	log.Printf("Authorized on acc %v", botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[This is log for msg] [%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text+update.Message.From.UserName)
		msg.ReplyToMessageID = update.Message.MessageID

		botAPI.Send(msg)
	}
}
