package config

import (
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	singletonInstance = &Config{}
	once              = sync.Once{}
)

type Config struct {
	Telegram
	Database
	NotifyWorker
}

type NotifyWorker struct {
	Period time.Duration `env:"NOTIFY_WORKER_PERIOD"`
}

type Telegram struct {
	AccessToken      string `env:"TELEGRAM_ACCESS_TOKEN"`
	UpdatesIntervals int    `env:"TELEGRAM_ACCESS_INTERVALS" env_default:"60"`
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
