package scheduler

import (
	mathbattle "mathbattle/models"
)

type MessageScheduler struct {
	repository mathbattle.ScheduledMessageRepository
}

func NewMessageScheduler(repository mathbattle.ScheduledMessageRepository) MessageScheduler {
	return MessageScheduler{
		repository: repository,
	}
}

func (s *MessageScheduler) Schedule(message mathbattle.ScheduledMessage) error {
	return nil
}

func (s *MessageScheduler) Start() error {
	return nil
}
