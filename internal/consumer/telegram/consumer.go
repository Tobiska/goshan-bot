package telegram

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"goshan-bot/internal/models"
)

const (
	getUpdatesInterval = 500 * time.Millisecond
)

type telegramClient interface {
	GetUpdates(context.Context, int) ([]tgbotapi.Update, error)
}

type router interface {
	Route(context.Context, *models.IncomingMessage)
}

type Consumer struct {
	telegramClient telegramClient
	router         router
}

func New(telegramClient telegramClient, router router) *Consumer {
	return &Consumer{
		telegramClient: telegramClient,
		router:         router,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	currentOffset := 0
	for ctx.Err() == nil {
		updates, err := c.telegramClient.GetUpdates(ctx, currentOffset+1)
		if err != nil {
			return fmt.Errorf("get updates from telegram error: %w", err)
		}

		for _, upd := range updates {
			if upd.Message != nil {
				c.router.Route(ctx, &models.IncomingMessage{
					ChatID:          upd.Message.Chat.ID,
					Username:        upd.Message.From.UserName,
					Text:            upd.Message.Text,
					UsernameDisplay: upd.Message.From.UserName,
					UserID:          upd.Message.From.ID,
				})
			}

			if upd.CallbackQuery != nil {
				panic("doesn't implemented")
			}

			currentOffset = upd.UpdateID
		}
		c.wait()
	}
	return nil
}

func (c *Consumer) wait() {
	time.Sleep(getUpdatesInterval)
}
