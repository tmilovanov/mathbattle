package mem

import (
	mathbattle "mathbattle/models"
)

type UserContextRepository struct {
	userContexts map[int64]mathbattle.TelegramUserContext
}

func NewUserContextRepository() (UserContextRepository, error) {
	return UserContextRepository{
		userContexts: make(map[int64]mathbattle.TelegramUserContext),
	}, nil
}

func (r *UserContextRepository) GetByTelegramID(chatID int64) (mathbattle.TelegramUserContext, error) {
	if ctx, isExist := r.userContexts[chatID]; isExist {
		return ctx, nil
	}

	newCtx := mathbattle.TelegramUserContext{
		ChatID:         chatID,
		Variables:      make(map[string]mathbattle.ContextVariable),
		CurrentStep:    0,
		CurrentCommand: "",
	}
	r.userContexts[chatID] = newCtx

	return newCtx, nil
}

func (r *UserContextRepository) Update(chatID int64, ctx mathbattle.TelegramUserContext) error {
	r.userContexts[chatID] = ctx
	return nil
}
