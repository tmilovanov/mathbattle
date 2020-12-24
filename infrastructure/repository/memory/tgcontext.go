package memory

import (
	"mathbattle/infrastructure"

	"mathbattle/models/mathbattle"
)

type TelegramContextRepository struct {
	userContexts   map[int64]infrastructure.TelegramUserContext
	userRepository mathbattle.UserRepository
}

func NewTelegramContextRepository(userRepository mathbattle.UserRepository) (TelegramContextRepository, error) {
	return TelegramContextRepository{
		userContexts:   make(map[int64]infrastructure.TelegramUserContext),
		userRepository: userRepository,
	}, nil
}

func (r *TelegramContextRepository) GetByTelegramID(chatID int64) (infrastructure.TelegramUserContext, error) {
	if ctx, isExist := r.userContexts[chatID]; isExist {
		return ctx, nil
	}

	user, err := r.userRepository.GetOrCreateByTelegramID(chatID)
	if err != nil {
		return infrastructure.TelegramUserContext{}, err
	}

	newCtx := infrastructure.TelegramUserContext{
		User:           user,
		Variables:      make(map[string]infrastructure.ContextVariable),
		CurrentStep:    0,
		CurrentCommand: "",
	}
	r.userContexts[chatID] = newCtx

	return newCtx, nil
}

func (r *TelegramContextRepository) Update(chatID int64, ctx infrastructure.TelegramUserContext) error {
	r.userContexts[chatID] = ctx
	return nil
}
