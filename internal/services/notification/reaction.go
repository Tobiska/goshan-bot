package notification

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"goshan-bot/internal/models"
)

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

func (s *Service) reactOnDescriptionSuccess(ctx context.Context, message models.IncomingMessage) error {
	if err := s.telegramSender.SendTextMessage(ctx, message.ChatID, AddEventAtNextStepMessage); err != nil {
		return fmt.Errorf("error while send message: %w", err)
	}
	return nil
}

func (s *Service) reactOnEventAtSuccess(ctx context.Context, message models.IncomingMessage) error {
	msgConfig := tgbotapi.NewMessage(message.ChatID, AddRemindInNextStepMessage)
	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1h", "1h"),
			tgbotapi.NewInlineKeyboardButtonData("30m", "1h"),
			tgbotapi.NewInlineKeyboardButtonData("5m", "1h"),
		),
	)
	if err := s.telegramSender.SendMessage(ctx, msgConfig); err != nil {
		return fmt.Errorf("error while send message: %w", err)
	}
	return nil
}

func (s *Service) reactOnRemindMeSuccess(ctx context.Context, message models.IncomingMessage) error {
	if err := s.telegramSender.SendTextMessage(ctx, message.ChatID, BuildDoneMessage); err != nil {
		return fmt.Errorf("error while send message: %w", err)
	}
	return nil
}

func (s *Service) reactOnTag(ctx context.Context, message models.IncomingMessage) error {
	if err := s.telegramSender.SendTextMessage(ctx, message.ChatID, AddDescriptionNextStepMessage); err != nil {
		return fmt.Errorf("error while send message: %w", err)
	}

	return nil
}
