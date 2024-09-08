package user

import (
	"context"
	"fmt"

	"goshan-bot/internal/models"
)

const (
	UserAlreadyExistMessage     = "Кажется вы уже были зарегистрированы ранее. Попробуйте создать нотификацию с помощью /add."
	UserSaveSuccessfullyMessage = "Привет! Я Гошан-Бот - помогаю не забывать о важных встречах или событиях. Попробуй создать нотификацию с помощью /add."
	UserSaveFailedMessage       = "Кажется, что-то пошло не так! Я пытаюсь это починить! Вернись попозже :D"
)

type userRepository interface {
	Save(context.Context, models.User) error
	FindByID(context.Context, int64) (*models.User, error)
}

type telegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type Service struct {
	userRepository userRepository
	telegramSender telegramSender
}

func New(userRepository userRepository, telegramSender telegramSender) *Service {
	return &Service{
		userRepository: userRepository,
		telegramSender: telegramSender,
	}
}

func (s *Service) StartCommand(ctx context.Context, message models.IncomingMessage) error {
	var err error

	defer func() {
		if err != nil {
			_ = s.telegramSender.SendMessage(ctx, message.ChatID, UserSaveFailedMessage)
		}
	}()

	u, err := s.userRepository.FindByID(ctx, message.UserID)
	if err != nil {
		return fmt.Errorf("user find by id error: %w", err)
	}

	if u == nil {
		if err := s.userRepository.Save(ctx, models.User{
			ID:       message.UserID,
			ChatID:   message.ChatID,
			Username: message.Username,
		}); err != nil {
			return fmt.Errorf("user save error: %w", err)
		}

		if err := s.telegramSender.SendMessage(ctx, u.ChatID, UserSaveSuccessfullyMessage); err != nil {
			return fmt.Errorf("user send message error: %w", err)
		}
		return nil

	}

	if err := s.telegramSender.SendMessage(ctx, u.ChatID, UserAlreadyExistMessage); err != nil {
		return fmt.Errorf("user send message error: %w", err)
	}

	return nil
}
