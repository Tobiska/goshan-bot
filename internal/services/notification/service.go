package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	"goshan-bot/internal/models"
)

// /add -> tag ->

const (
	BuildAlreadyExistMessage   = "Вы уже находитесь в процессе добавления новой нотификации, продолжайте заполнять недостающую информацию :D"
	BuildNotExistMessage       = "Я пока не умею отвечать на обычные сообщения и вести диалог выполни одну из доступных команд."
	EventAtParseErrorAtMessage = "Кажется ты ввёл дату не в правильном формате :("
	RemindMeParseErrorMessage  = "Кажется ты ввёл продолжительность в неверном формате :((("
)

const (
	AddTagNextStepMessage         = "Введи уникальный тег для вашего уведомления."
	AddDescriptionNextStepMessage = "Введи описание или примечание для вашего уведомления. Это может быть например ссылка на онлайн встречу или номер аудитории."
	AddEventAtNextStepMessage     = "Введи время встречи или события в формате. 2006-01-02 15:04:05 по НСК."
	AddRemindInNextStepMessage    = "Напиши за какое время тебя предупредить о событии. Формат: 2h5m1s."
	BuildDoneMessage              = "Отлично вся нужная мне информация собрана!"
)

var (
	ErrAlreadyExist = errors.New("already exist")
)

type telegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type userRepository interface {
	FindByID(context.Context, int64) (*models.User, error)
}

type notificationRepository interface {
	CreateBuild(ctx context.Context, builder models.NotificationBuild) error
	UpdateBuild(ctx context.Context, builder models.NotificationBuild) error
	CreateNotification(ctx context.Context, notification *models.Notification) error
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

	if err != nil {
		return fmt.Errorf("error while add create build notification: %w", err)
	}
	if errors.Is(err, ErrAlreadyExist) {
		if err := s.telegramSender.SendMessage(ctx, message.ChatID, BuildAlreadyExistMessage); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	}

	if err := s.telegramSender.SendMessage(ctx, message.ChatID, AddTagNextStepMessage); err != nil {
		return fmt.Errorf("error while send message: %w", err)
	}

	return nil
}

func (s *Service) HandleMessage(ctx context.Context, message models.IncomingMessage) error {
	build, err := s.notificationRepository.FindBuildByUserID(ctx, message.UserID)
	if err != nil {
		return fmt.Errorf("error while find build by userID: %w", err)
	}

	if build == nil {
		if err := s.telegramSender.SendMessage(ctx, message.ChatID, BuildNotExistMessage); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
		return nil
	}

	if build.Tag != nil {
		err = s.notificationRepository.UpdateBuild(ctx, models.NotificationBuild{
			Tag: &message.Text,
		})
		if err != nil {
			return fmt.Errorf("error while update build: %w", err)
		}

		if err := s.telegramSender.SendMessage(ctx, message.ChatID, AddDescriptionNextStepMessage); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	} else if build.Description != nil {
		err = s.notificationRepository.UpdateBuild(ctx, models.NotificationBuild{
			Description: &message.Text,
		})
		if err != nil {
			return fmt.Errorf("error while update build: %w", err)
		}

		if err := s.telegramSender.SendMessage(ctx, message.ChatID, AddEventAtNextStepMessage); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	} else if build.EventAt != nil {
		eventAt, err := time.Parse(time.DateTime, message.Text)
		if err != nil {
			if err := s.telegramSender.SendMessage(ctx, message.ChatID, EventAtParseErrorAtMessage); err != nil {
				return fmt.Errorf("error while send message: %w", err)
			}
			return nil
		}

		err = s.notificationRepository.UpdateBuild(ctx, models.NotificationBuild{
			EventAt: &eventAt,
		})
		if err != nil {
			return fmt.Errorf("error while update build: %w", err)
		}

		if err := s.telegramSender.SendMessage(ctx, message.ChatID, AddRemindInNextStepMessage); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	} else if build.RemindIn != nil {
		remindMe, err := time.ParseDuration(message.Text)
		if err != nil {
			if err := s.telegramSender.SendMessage(ctx, message.ChatID, RemindMeParseErrorMessage); err != nil {
				return fmt.Errorf("error while send message: %w", err)
			}
			return nil
		}

		build.RemindIn = &remindMe

		err = s.notificationRepository.CreateNotification(ctx, models.CreateFromBuild(build))
		if err != nil {
			return fmt.Errorf("create notification error: %w", err)
		}

		if err := s.telegramSender.SendMessage(ctx, message.ChatID, BuildDoneMessage); err != nil {
			return fmt.Errorf("error while send message: %w", err)
		}
	}

	return nil
}
