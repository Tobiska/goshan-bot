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

func (c *Client) SendTextMessage(_ context.Context, chatID int64, text string) error {
	msgConfig := tgbotapi.NewMessage(chatID, text)

	_, err := c.bot.Send(msgConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RequestCallback(_ context.Context, id, text string) error {
	callback := tgbotapi.NewCallback(id, text)
	if _, err := c.bot.Request(callback); err != nil {
		return fmt.Errorf("error while request telegram api: %w", err)
	}
	return nil
}

func (c *Client) DeleteInlineButtons(_ context.Context, chatID int64, messageID int) error {
	editMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{}))
	if _, err := c.bot.Send(editMsg); err != nil {
		return fmt.Errorf("send edit message error: %w", err)
	}

	return nil
}

func (c *Client) SendMessage(_ context.Context, message tgbotapi.MessageConfig) error {
	_, err := c.bot.Send(message)
	if err != nil {
		return err
	}
	return nil
}
