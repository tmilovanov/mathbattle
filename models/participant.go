package models

import (
	"strconv"
	"time"
	"unicode"
)

type Participant struct {
	ID               string
	TelegramID       int64
	Name             string
	School           string
	Grade            int
	RegistrationTime time.Time
}

type ParticipantRepository interface {
	Store(participant Participant) (Participant, error)
	GetByID(ID string) (Participant, error)
	GetByTelegramID(TelegramID int64) (Participant, error)
	GetAll() ([]Participant, error)
	Update(participant Participant) error
	Delete(ID string) error
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

func IsRegistered(participantRepository ParticipantRepository, telegramID int64) (bool, error) {
	_, err := participantRepository.GetByTelegramID(telegramID)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func FilterRegisteredAfter(participants []Participant, datetime time.Time) []Participant {
	result := []Participant{}

	for _, participant := range participants {
		if participant.RegistrationTime.After(datetime) {
			result = append(result, participant)
		}
	}

	return result
}
