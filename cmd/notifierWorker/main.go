package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"goshan-bot/internal/clients/sqlite"
	telegramClient "goshan-bot/internal/clients/telegram"
	"goshan-bot/internal/config"
	notificationRepository "goshan-bot/internal/repository/notification"
	userRepository "goshan-bot/internal/repository/user"
	notificationService "goshan-bot/internal/services/notification"
	"goshan-bot/internal/sheduler"
)

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer cancel()

	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("read while read config: %w", err)
	}

	tgCli, err := telegramClient.New(&cfg.Telegram)
	if err != nil {
		return fmt.Errorf("error while initialize telegram client: %w", err)
	}

	db, err := sqlite.New(&cfg.Database)
	if err != nil {
		return fmt.Errorf("error initialize: %w", err)
	}

	userRepo := userRepository.New(db)

	notificationRepo := notificationRepository.New(db)

	notificationSrv := notificationService.New(tgCli, userRepo, notificationRepo)

	err = sheduler.New(cfg.Period).Run(ctx, func(ctx context.Context) {
		if err := notificationSrv.Notify(ctx); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return fmt.Errorf("error while run scheduler: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
