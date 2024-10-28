package main

import (
	"fmt"
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
		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "Бот запущен"
		case "check_id":
			userID := update.Message.From.ID
			log.Printf("[This is log for msg] [%s] %s\n", update.Message.From.UserName, update.Message.Text)
			msg.Text = fmt.Sprintf("Ваш ID: %v", userID)
			log.Printf("[This is log for messageBot]: %v\n", msg)
		default:
			msg.Text = "Я не знаю такой команды"
		}

		if _, err := botAPI.Send(msg); err != nil {
			log.Printf("[ERROR] Failed to send ID: %v", err)
		}
	}
}
