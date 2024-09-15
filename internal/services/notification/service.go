package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"goshan-bot/internal/models"
)

var (
	ErrAlreadyExist = errors.New("already exist")
)

type telegramSender interface {
	SendTextMessage(ctx context.Context, chatID int64, text string) error
	SendMessage(ctx context.Context, message tgbotapi.MessageConfig) error
}

type userRepository interface {
	FindByID(context.Context, int64) (*models.User, error)
}

type notificationRepository interface {
	CreateBuild(ctx context.Context, builder models.NotificationBuild) error
	UpdateBuild(ctx context.Context, chatID, userID int64, builder models.NotificationBuild) error
	GetClosestNotification(ctx context.Context, limit, offset int64, now time.Time) ([]models.Notification, error)
	CreateNotification(ctx context.Context, notification *models.Notification) error
	DeleteBuild(ctx context.Context, chatID, userID int64) error
	Delete(ctx context.Context, chatID, userID int64) error
	FindBuildByUserID(ctx context.Context, userID int64) (*models.NotificationBuild, error)
}

type Service struct {
	telegramSender         telegramSender
	userRepository         userRepository
	notificationRepository notificationRepository
}

func New(telegramSender telegramSender, userRepository userRepository, notificationRepository notificationRepository) *Service {
	return &Service{
		telegramSender:         telegramSender,
		userRepository:         userRepository,
		notificationRepository: notificationRepository,
	}
}

func (s *Service) AddCommand(ctx context.Context, message models.IncomingMessage) error {
	var err error

	err = s.notificationRepository.CreateBuild(ctx, models.NotificationBuild{
		UserID: message.UserID,
		ChatID: message.ChatID,
	})

	if errors.Is(err, ErrAlreadyExist) {
		if err := s.telegramSender.SendTextMessage(ctx, message.ChatID, BuildAlreadyExistMessage); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("error while add create build notification: %w", err)
	}

	if err := s.telegramSender.SendTextMessage(ctx, message.ChatID, AddTagNextStepMessage); err != nil {
		return fmt.Errorf("error while send message: %w", err)
	}

	return nil
}

func (s *Service) BuildNotification(ctx context.Context, message models.IncomingMessage) error {
	build, err := s.notificationRepository.FindBuildByUserID(ctx, message.UserID)
	if err != nil {
		return fmt.Errorf("error while find build by userID: %w", err)
	}

	if build == nil {
		if err := s.telegramSender.SendTextMessage(ctx, message.ChatID, BuildNotExistMessage); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
		return nil
	}

	if build.Tag == nil {
		err = s.notificationRepository.UpdateBuild(ctx, message.ChatID, message.UserID, models.NotificationBuild{
			Tag: &message.Text,
		})
		if err != nil {
			return fmt.Errorf("error while update build: %w", err)
		}

		if err := s.reactOnTag(ctx, message); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	} else if build.Description == nil {
		err = s.notificationRepository.UpdateBuild(ctx, message.ChatID, message.UserID, models.NotificationBuild{
			Description: &message.Text,
		})
		if err != nil {
			return fmt.Errorf("error while update build: %w", err)
		}

		if err := s.reactOnDescriptionSuccess(ctx, message); err != nil {
			return fmt.Errorf("error while react: %w", err)
		}
	} else if build.EventAt == nil {
		eventAt, err := time.Parse(time.DateTime, message.Text)
		if err != nil {
			if err := s.telegramSender.SendTextMessage(ctx, message.ChatID, EventAtParseErrorAtMessage); err != nil {
				return fmt.Errorf("error while send message: %w", err)
			}
			return nil
		}

		err = s.notificationRepository.UpdateBuild(ctx, message.ChatID, message.UserID, models.NotificationBuild{
			EventAt: &eventAt,
		})
		if err != nil {
			return fmt.Errorf("error while update build: %w", err)
		}

		if err := s.reactOnEventAtSuccess(ctx, message); err != nil {
			return fmt.Errorf("error while react: %w", err)
		}
	} else if build.RemindIn == nil {
		remindMe, err := time.ParseDuration(message.Text)
		if err != nil {
			if err := s.telegramSender.SendTextMessage(ctx, message.ChatID, RemindMeParseErrorMessage); err != nil {
				return fmt.Errorf("error while send message: %w", err)
			}
			return nil
		}

		build.RemindIn = &remindMe

		err = s.notificationRepository.CreateNotification(ctx, models.CreateFromBuild(build))
		if err != nil {
			return fmt.Errorf("create notification error: %w", err)
		}

		err = s.notificationRepository.DeleteBuild(ctx, message.ChatID, message.UserID)
		if err != nil {
			return fmt.Errorf("create notification error: %w", err)
		}

		if err := s.reactOnRemindMeSuccess(ctx, message); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	}

	return nil
}

func (s *Service) Notify(ctx context.Context) error {
	now := time.Now()

	notifications, err := s.notificationRepository.GetClosestNotification(ctx, 20, 0, now)
	if err != nil {
		return fmt.Errorf("error while get closest notifications: %w", err)
	}

	for _, n := range notifications {
		if err := s.telegramSender.SendTextMessage(ctx, n.ChatID, buildNotifyMessage(n, now)); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}

		if err := s.notificationRepository.Delete(ctx, n.ChatID, n.UserID); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	}

	return nil
}

func buildNotifyMessage(n models.Notification, now time.Time) string {
	return fmt.Sprintf(`
		Эй! Не забудь что через %s будет %s - %s
    `, n.EventAt.Sub(now).Abs().String(), n.Tag, n.Description)
}
