package config

import (
	"log"
	"sync"

	"github.com/cristalhq/aconfig"
)

type Config struct {
	TelegramBotToken string `default:"telegram_bot_token" usage:"Token" required:"true"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			SkipEnv:   true,
			SkipFlags: true,
			Files:     []string{"config.json"},
		})
		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] failed to load config: %v", err)
		}
	})

	return cfg
}
