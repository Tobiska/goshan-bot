package router

import (
	"context"

	"goshan-bot/internal/models"
)

type telegramClient interface {
	SendMessage(context.Context, int64, string) error
}

type Router struct {
	telegramClient telegramClient
}

func New(telegramClient telegramClient) *Router {
	return &Router{
		telegramClient: telegramClient,
	}
}

func (r *Router) Route(ctx context.Context, msg *models.IncomingMessage) {
	if err := r.telegramClient.SendMessage(ctx, msg.ChatID, msg.Text); err != nil {
	}
}
