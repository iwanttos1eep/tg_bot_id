package config

import (
	"log"
	"sync"

	"github.com/cristalhq/aconfig"
)

type Config struct {
	TelegramBotToken string `default:"7784142806:AAHEOahcjkeZQm2hcmOq7Yq9MAuUfIAQVvQ" usage:"Token" required:"true"`
	//TelegramChannelID    	int64           `hcl:"telegram_channel_id" env:"TELEGRAM_CHANNEL_ID" required:"true"`
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
			Files:     []string{"./config.json"},
		})

		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] failed to load config: %v", err)
		}
	})

	return cfg
}
