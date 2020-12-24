package mathbattle

import (
	"strconv"
	"time"
	"unicode"
)

type Participant struct {
	ID               string    `json:"id"`
	TelegramID       int64     `json:"telegram_id"`
	Name             string    `json:"name"`
	School           string    `json:"school"`
	Grade            int       `json:"grade"`
	RegistrationTime time.Time `json:"registration_time"`
}

type ParticipantRepository interface {
	Store(participant Participant) (Participant, error)
	GetByID(ID string) (Participant, error)
	GetByTelegramID(TelegramID int64) (Participant, error)
	GetAll() ([]Participant, error)
	Update(participant Participant) error
	Delete(ID string) error
}

type ParticipantService interface {
	Store(participant Participant) (Participant, error)
	GetByID(ID string) (Participant, error)
	GetByTelegramID(TelegramID int64) (Participant, error)
	GetAll() ([]Participant, error)
	Update(participant Participant) error
	Delete(ID string) error
	Unsubscribe(ID string) error
}

func IsValidGrade(grade int) bool {
	if grade >= 1 && grade <= 11 {
		return true
	}
	return false
}

func IsParticipantNameValid(input string) bool {
	letters := []rune(input)

	if len(letters) == 0 || len(letters) > 30 {
		return false
	}

	for _, r := range []rune(input) {
		if !(unicode.IsLetter(r) || r == ' ') {
			return false
		}
	}
	return true
}

func ValidateUserName(userInput string) (string, bool) {
	if isOk := IsParticipantNameValid(userInput); !isOk {
		return "", false
	}
	return userInput, true
}

func ValidateUserGrade(userInput string) (int, bool) {
	r, err := strconv.Atoi(userInput)
	if err != nil {
		return 0, false
	}
	if isOk := IsValidGrade(r); !isOk {
		return 0, false
	}
	return r, true
}
