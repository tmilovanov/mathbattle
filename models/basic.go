package models

import (
	"errors"
	"strconv"
	"time"
	"unicode"
)

type Storage struct {
	Participants ParticipantRepository
	Rounds       RoundRepository
	Problems     ProblemRepository
	Solutions    SolutionRepository
}

type ParticipantRepository interface {
	Store(participant Participant) error
	GetByID(ID string) (Participant, bool, error)
	GetByTelegramID(TelegramID string) (Participant, bool, error)
	GetAll() ([]Participant, error)
	Delete(ID string) error
}

type ProblemRepository interface {
	Store(problem Problem) error
	GetByID(ID string) (Problem, error)
	GetAll() ([]Problem, error)
}

var ErrSolutionNotFound = errors.New("solution not found")

type SolutionRepository interface {
	Store(solution Solution) (Solution, error) // Return newly created Solution with filled in ID
	Get(ID string) (Solution, error)
	Find(roundID string, participantID string, problemID string) (Solution, error)
	FindOrCreate(roundID string, participantID string, problemID string) (Solution, error)
	AppendPart(ID string, part Image) error
	Delete(ID string) error
}

var ErrRoundNotFound = errors.New("round not found")

type RoundRepository interface {
	Store(round Round) error
	GetRunning() (Round, error)
	GetAll() ([]Round, error)
}

type Participant struct {
	ID               string
	TelegramID       string
	Name             string
	School           string
	Grade            int
	RegistrationTime time.Time
}

type Problem struct {
	ID        string
	MinGrade  int
	MaxGrade  int
	Extension string
	Content   []byte
}

type Image struct {
	Extension string
	Content   []byte
}

type Solution struct {
	ID            string
	ParticipantID string
	ProblemID     string
	RoundID       string
	Parts         []Image
}

type Round struct {
	ID                  string
	StartDate           time.Time
	EndDate             time.Time
	ProblemDistribution RoundDistribution
}

func NewRound() Round {
	return Round{
		ID:                  "", // ID is filled by database layer
		StartDate:           time.Now(),
		ProblemDistribution: make(map[string][]string),
	}
}

// RoundDistribution is mapping from participant id to set of problems that were sent to him
type RoundDistribution map[string][]string

type ProblemDistributor interface {
	Get(participants []Participant, problems []Problem, rounds []Round) (RoundDistribution, error)
}

func IsProblemSuitableForParticipant(problem *Problem, participant *Participant) bool {
	if participant.Grade >= problem.MinGrade && participant.Grade <= problem.MaxGrade {
		return true
	}
	return false
}

func IsParticipantNameValid(input string) bool {
	for _, r := range []rune(input) {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsValidGrade(grade int) bool {
	if grade >= 1 && grade <= 11 {
		return true
	}
	return false
}

func GetProblemIDs(problems []Problem) []string {
	result := []string{}
	for _, problem := range problems {
		result = append(result, problem.ID)
	}
	return result
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
	_, exist, err := participantRepository.GetByTelegramID(strconv.FormatInt(telegramID, 10))
	if err != nil {
		return false, err
	}

	return exist, nil
}
