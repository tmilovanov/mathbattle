package application

import "mathbattle/models/mathbattle"

type ParticipantService struct {
	Rep mathbattle.ParticipantRepository
}

func (ps *ParticipantService) Store(participant mathbattle.Participant) (mathbattle.Participant, error) {
	return ps.Rep.Store(participant)
}

func (ps *ParticipantService) GetByID(ID string) (mathbattle.Participant, error) {
	return ps.Rep.GetByID(ID)
}

func (ps *ParticipantService) GetByTelegramID(TelegramID int64) (mathbattle.Participant, error) {
	return ps.Rep.GetByTelegramID(TelegramID)
}

func (ps *ParticipantService) GetAll() ([]mathbattle.Participant, error) {
	return ps.Rep.GetAll()
}

func (ps *ParticipantService) Update(participant mathbattle.Participant) error {
	return ps.Rep.Update(participant)
}

func (ps *ParticipantService) Delete(ID string) error {
	return ps.Rep.Delete(ID)
}

func (ps *ParticipantService) Unsubscribe(ID string) error {
	return ps.Rep.Delete(ID)
}
