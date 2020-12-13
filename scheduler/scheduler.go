package scheduler

import (
	"log"
	mathbattle "mathbattle/models"
	"time"
)

type MessageScheduler struct {
	repository   mathbattle.ScheduledMessageRepository
	participants mathbattle.ParticipantRepository
	postman      mathbattle.TelegramPostman
}

func NewMessageScheduler(repository mathbattle.ScheduledMessageRepository, participants mathbattle.ParticipantRepository,
	postman mathbattle.TelegramPostman) MessageScheduler {
	return MessageScheduler{
		repository:   repository,
		participants: participants,
		postman:      postman,
	}
}

func (s *MessageScheduler) scheduleSend(message mathbattle.ScheduledMessage) {
	time.AfterFunc(time.Until(message.SendTime), func(msg mathbattle.ScheduledMessage) func() {
		return func() {
			participants, err := s.participants.GetAll()
			if err != nil {
				log.Printf("Failed to send scheduled message, error: %v", err)
			}
			for _, participant := range participants {
				s.postman.PostText(participant.TelegramID, msg.Message)
			}
		}
	}(message))
}

func (s *MessageScheduler) Schedule(message mathbattle.ScheduledMessage) error {
	log.Println("Schedule()")
	msg, err := s.repository.Store(message)
	if err != nil {
		return err
	}

	s.scheduleSend(msg)
	return nil
}

func (s *MessageScheduler) StartAll() error {
	messages, err := s.repository.GetAll()
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if msg.SendTime.Before(time.Now()) {
			log.Printf("No need to run. Now: %v, SendTIme: %v", time.Now(), msg.SendTime)
		} else {
			s.scheduleSend(msg)
		}
	}
	return nil
}
