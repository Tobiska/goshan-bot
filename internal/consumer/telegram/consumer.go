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
	DeleteInlineButtons(_ context.Context, chatID int64, messageID int) error
}

type router interface {
	RouteMessage(context.Context, *models.IncomingMessage)
	RouteCallback(ctx context.Context, msg *models.IncomingMessage)
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
				c.router.RouteMessage(ctx, &models.IncomingMessage{
					ChatID:          upd.Message.Chat.ID,
					Username:        upd.Message.From.UserName,
					Text:            upd.Message.Text,
					UsernameDisplay: upd.Message.From.UserName,
					UserID:          upd.Message.From.ID,
					IsCallback:      false,
				})
			}

			if upd.CallbackQuery != nil {
				if err := c.telegramClient.DeleteInlineButtons(ctx, upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID); err != nil {
					return fmt.Errorf("error while deleting inline buttons from message: %w", err)
				}

				c.router.RouteMessage(ctx, &models.IncomingMessage{
					ChatID:          upd.CallbackQuery.Message.Chat.ID,
					UserID:          upd.CallbackQuery.Message.From.ID,
					Username:        upd.CallbackQuery.From.UserName,
					UsernameDisplay: upd.CallbackQuery.From.UserName,
					Text:            upd.CallbackQuery.Data,
					IsCallback:      true,
					CallbackMsgID:   upd.CallbackQuery.Message.MessageID,
				})

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
