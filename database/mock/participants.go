package mock

import (
	mathbattle "mathbattle/models"

	"github.com/pkg/errors"
)

type MockParticipantsRepository struct {
	impl map[string]mathbattle.Participant
}

func NewMockParticipantsRepository() MockParticipantsRepository {
	return MockParticipantsRepository{make(map[string]mathbattle.Participant)}
}

func (r *MockParticipantsRepository) Store(participant mathbattle.Participant) (mathbattle.Participant, error) {
	toStore := participant
	toStore.ID = participant.TelegramID
	r.impl[participant.TelegramID] = toStore
	return toStore, nil
}

func (r *MockParticipantsRepository) GetByID(ID string) (mathbattle.Participant, error) {
	p, ok := r.impl[ID]
	if !ok {
		return p, mathbattle.ErrNotFound
	}
	return p, nil
}

func (r *MockParticipantsRepository) GetByTelegramID(TelegramID string) (mathbattle.Participant, error) {
	for item := range r.impl {
		if r.impl[item].TelegramID == TelegramID {
			return r.impl[item], nil
		}
	}

	return mathbattle.Participant{}, mathbattle.ErrNotFound
}

func (r *MockParticipantsRepository) GetAll() ([]mathbattle.Participant, error) {
	return []mathbattle.Participant{}, errors.Errorf("Not implemented")
}

func (r *MockParticipantsRepository) Delete(ID string) error {
	return errors.Errorf("Not implemented")
}
