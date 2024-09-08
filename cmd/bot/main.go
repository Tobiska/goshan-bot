package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	telegramClient "goshan-bot/internal/clients/telegram"
	"goshan-bot/internal/config"
	telegramConsumer "goshan-bot/internal/consumer/telegram"
	"goshan-bot/internal/router"
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

	rt := router.New(tgCli)

	cs := telegramConsumer.New(tgCli, rt)

	if err := cs.Run(ctx); err != nil {
		return fmt.Errorf("consumer run error: %w", err)
	}

	return nil
}
