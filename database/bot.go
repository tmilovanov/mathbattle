package database

import (
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type InMemoryRepository struct {
	userContexts map[int64]mathbattle.TelegramUserContext
}

func NewInMemoryRepository() (InMemoryRepository, error) {
	return InMemoryRepository{
		userContexts: make(map[int64]mathbattle.TelegramUserContext),
	}, nil
}

func (r *InMemoryRepository) GetByTelegramID(chatID int64, bot *tb.Bot) (mathbattle.TelegramUserContext, error) {
	if ctx, isExist := r.userContexts[chatID]; isExist {
		return ctx, nil
	}

	newCtx := mathbattle.TelegramUserContext{
		ChatID:         chatID,
		Variables:      make(map[string]mathbattle.ContextVariable),
		CurrentStep:    0,
		CurrentCommand: "",
		Bot:            bot,
	}
	r.userContexts[chatID] = newCtx

	return newCtx, nil
}

func (r *InMemoryRepository) Update(ctx mathbattle.TelegramUserContext) error {
	r.userContexts[ctx.ChatID] = ctx
	return nil
}
