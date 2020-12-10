package sqlite

import (
	mathbattle "mathbattle/models"
)

type ScheduledMessageRepository struct {
	sqliteRepository
}

func NewScheduledMessageRepository(dbPath string) (ScheduledMessageRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return ScheduledMessageRepository{}, err
	}

	return ScheduledMessageRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func (mr *ScheduledMessageRepository) Store(msg mathbattle.ScheduledMessage) (mathbattle.ScheduledMessage, error) {
	return mathbattle.ScheduledMessage{}, nil
}

func (mr *ScheduledMessageRepository) Get(ID string) (mathbattle.ScheduledMessage, error) {
	return mathbattle.ScheduledMessage{}, nil
}

func (mr *ScheduledMessageRepository) Update(msg mathbattle.ScheduledMessage) error {
	return nil
}

func (mr *ScheduledMessageRepository) Delete(ID string) error {
	return nil
}
