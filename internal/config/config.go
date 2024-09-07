package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	singletonInstance = &Config{}
	once              = sync.Once{}
)

type Config struct {
	Telegram
}

type Telegram struct {
	AccessToken      string `env:"TELEGRAM_ACCESS_TOKEN"`
	UpdatesIntervals int    `env:"TELEGRAM_ACCESS_INTERNALS"`
}

func ReadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		err = cleanenv.ReadEnv(singletonInstance)
	})

	return singletonInstance, err
}
