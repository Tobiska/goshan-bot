package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"goshan-bot/internal/config"
)

type Client struct {
	bot *tgbotapi.BotAPI
	cfg *config.Telegram
}

func New(cfg *config.Telegram) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("initialize bot error: %w", err)
	}
	bot.Debug = cfg.Debug

	return &Client{
		bot: bot,
		cfg: cfg,
	}, nil
}

func (c *Client) GetUpdates(_ context.Context, offset int) ([]tgbotapi.Update, error) {
	updateConfig := tgbotapi.NewUpdate(offset) // update offset
	updateConfig.Timeout = c.cfg.UpdatesIntervals

	return c.bot.GetUpdates(updateConfig)
}

func (c *Client) SendMessage(_ context.Context, chatID int64, text string) error {
	msgConfig := tgbotapi.NewMessage(chatID, text)

	_, err := c.bot.Send(msgConfig)
	if err != nil {
		return err
	}

	return nil
}
