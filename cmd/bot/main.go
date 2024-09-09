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
	telegramConsumer "goshan-bot/internal/consumer/telegram"
	notificationRepository "goshan-bot/internal/repository/notification"
	userRepository "goshan-bot/internal/repository/user"
	"goshan-bot/internal/router"
	notificationService "goshan-bot/internal/services/notification"
	userService "goshan-bot/internal/services/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

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

	userSrv := userService.New(userRepo, tgCli)

	rt := router.New(tgCli, userSrv, notificationSrv)

	cs := telegramConsumer.New(tgCli, rt)

	if err := cs.Run(ctx); err != nil {
		return fmt.Errorf("consumer run error: %w", err)
	}

	return nil
}
