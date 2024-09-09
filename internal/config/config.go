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
	Database
}

type Telegram struct {
	AccessToken      string `env:"TELEGRAM_ACCESS_TOKEN"`
	UpdatesIntervals int    `env:"TELEGRAM_ACCESS_INTERNALS"`
	Debug            bool   `env:"TELEGRAM_DEBUG"`
}

type Database struct {
	Dsn string `env:"DATABASE_DSN"`
}

func ReadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		err = cleanenv.ReadEnv(singletonInstance)
	})

	return singletonInstance, err
}
