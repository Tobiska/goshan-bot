package router

import (
	"context"
	"log"

	"goshan-bot/internal/models"
)

const (
	addCommand     = "/add"
	startCommand   = "/start"
	currentCommand = "/current"
	listCommand    = "/list"
	updateCommand  = "/update"
	deleteCommand  = "/delete"
)

type telegramClient interface {
	SendMessage(context.Context, int64, string) error
}

type userService interface {
	StartCommand(ctx context.Context, message models.IncomingMessage) error
}

type notificationService interface {
	AddCommand(ctx context.Context, message models.IncomingMessage) error
	HandleMessage(ctx context.Context, message models.IncomingMessage) error
}

type Router struct {
	telegramClient      telegramClient
	userService         userService
	notificationService notificationService
}

func New(telegramClient telegramClient, userService userService, notificationService notificationService) *Router {
	return &Router{
		telegramClient:      telegramClient,
		userService:         userService,
		notificationService: notificationService,
	}
}

func (r *Router) Route(ctx context.Context, msg *models.IncomingMessage) {
	if msg.Text == "/start" {
		if err := r.userService.StartCommand(ctx, *msg); err != nil {
			log.Println(err)
		}
		return
	}

	if msg.Text == "/add" {
		if err := r.notificationService.AddCommand(ctx, *msg); err != nil {
			log.Println(err)
		}
		return
	}

	if msg.Text != "" {
		if err := r.notificationService.HandleMessage(ctx, *msg); err != nil {
			log.Println(err)
		}
		return
	}
}
