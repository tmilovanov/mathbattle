package models

import "time"

type RecieversType string

const (
	Everyone RecieversType = "Everyone"
	Selected RecieversType = "Selected"
)

type ScheduledMessage struct {
	ID            string
	Message       string // Perfectly if message could be Text/photo/photoalbum/video/audio
	SendTime      time.Time
	RecieversType RecieversType
	Recievers     []string // Participant IDs if RecieversType == Selected
}

type ScheduledMessageRepository interface {
	Store(msg ScheduledMessage) (ScheduledMessage, error)
	Get(ID string) (ScheduledMessage, error)
	GetAll() ([]ScheduledMessage, error)
	Update(msg ScheduledMessage) error
	Delete(ID string) error
}

type MessageScheduler interface {
	Schedule(message ScheduledMessage) error
	StartAll() error
}
