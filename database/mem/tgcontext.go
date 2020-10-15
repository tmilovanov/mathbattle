package mem

import (
	mathbattle "mathbattle/models"
)

type TelegramContextRepository struct {
	userContexts   map[int64]mathbattle.TelegramUserContext
	userRepository mathbattle.TelegramUserRepository
}

func NewTelegramContextRepository(userRepository mathbattle.TelegramUserRepository) (TelegramContextRepository, error) {
	return TelegramContextRepository{
		userContexts:   make(map[int64]mathbattle.TelegramUserContext),
		userRepository: userRepository,
	}, nil
}

func (r *TelegramContextRepository) GetByTelegramID(chatID int64) (mathbattle.TelegramUserContext, error) {
	if ctx, isExist := r.userContexts[chatID]; isExist {
		return ctx, nil
	}

	user, err := r.userRepository.GetOrCreateByTelegramID(chatID)
	if err != nil {
		return mathbattle.TelegramUserContext{}, err
	}

	newCtx := mathbattle.TelegramUserContext{
		User:           user,
		Variables:      make(map[string]mathbattle.ContextVariable),
		CurrentStep:    0,
		CurrentCommand: "",
	}
	r.userContexts[chatID] = newCtx

	return newCtx, nil
}

func (r *TelegramContextRepository) Update(chatID int64, ctx mathbattle.TelegramUserContext) error {
	r.userContexts[chatID] = ctx
	return nil
}
