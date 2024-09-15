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
	SendTextMessage(context.Context, int64, string) error
}

type userService interface {
	StartCommand(ctx context.Context, message models.IncomingMessage) error
}

type notificationService interface {
	AddCommand(ctx context.Context, message models.IncomingMessage) error
	BuildNotification(ctx context.Context, message models.IncomingMessage) error
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

func (r *Router) RouteMessage(ctx context.Context, msg *models.IncomingMessage) {
	if msg.Text == startCommand {
		if err := r.userService.StartCommand(ctx, *msg); err != nil {
			log.Println(err)
		}
		return
	}

	if msg.Text == addCommand {
		if err := r.notificationService.AddCommand(ctx, *msg); err != nil {
			log.Println(err)
		}
		return
	}

	if msg.Text != "" {
		if err := r.notificationService.BuildNotification(ctx, *msg); err != nil {
			log.Println(err)
		}
		return
	}
}

func (r *Router) RouteCallback(ctx context.Context, msg *models.IncomingMessage) {
	if err := r.notificationService.BuildNotification(ctx, *msg); err != nil {
		log.Println(err)
	}
	return
}
