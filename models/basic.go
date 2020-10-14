package models

import (
	"errors"
	"strconv"
	"time"
	"unicode"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage struct {
	Participants ParticipantRepository
	Problems     ProblemRepository
	Solutions    SolutionRepository
	Rounds       RoundRepository
}

type ParticipantRepository interface {
	Store(participant Participant) (Participant, error)
	GetByID(ID string) (Participant, error)
	GetByTelegramID(TelegramID string) (Participant, error)
	GetAll() ([]Participant, error)
	Delete(ID string) error
}

type ProblemRepository interface {
	Store(problem Problem) (Problem, error)
	GetByID(ID string) (Problem, error)
	GetAll() ([]Problem, error)
}

type SolutionRepository interface {
	Store(solution Solution) (Solution, error) // Return newly created Solution with filled in ID
	Get(ID string) (Solution, error)
	Find(roundID string, participantID string, problemID string) (Solution, error)
	FindMany(roundID string, participantID string, problemID string) ([]Solution, error) //Leave IDs empty if it's not important
	FindOrCreate(roundID string, participantID string, problemID string) (Solution, error)
	AppendPart(ID string, part Image) error
	Delete(ID string) error
}

type RoundRepository interface {
	Store(round Round) (Round, error)
	Get(ID string) (Round, error)
	GetRunning() (Round, error)
	GetAll() ([]Round, error)
	Delete(roundID string) error
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
	ReviewDistribution  ReviewDistribution
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

type ReviewDistribution struct {
	BetweenParticipants map[string][]string // mapping from participant id to set of solutions that he is sent to review
	ToOrganizers        []Solution
}

type ProblemDistributor interface {
	Get(participants []Participant, problems []Problem, rounds []Round) (RoundDistribution, error)
}

type OnReviewDistributor interface {
	Get(participants []Participant, solutions []Solution, reviewersMinCount int) (ReviewDistribution, error)
}

// SolutionReview is review on solution that one participant sends to another
type SolutionReview struct {
	ID                    string
	ReviewerID            string
	ReviewedParticipantID string
	SolutionID            string
	Parts                 []Image
}

func IsProblemSuitableForParticipant(problem *Problem, participant *Participant) bool {
	if participant.Grade >= problem.MinGrade && participant.Grade <= problem.MaxGrade {
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

func ValidateProblemNumber(userInput string, problemIDs []string) (int, bool) {
	problemNumber, err := strconv.Atoi(userInput)
	if err != nil {
		return -1, false
	}
	problemNumber = problemNumber - 1
	if problemNumber < 0 || problemNumber >= len(problemIDs) {
		return -1, false
	}

	return problemNumber, true
}

func IsRegistered(participantRepository ParticipantRepository, telegramID int64) (bool, error) {
	_, err := participantRepository.GetByTelegramID(strconv.FormatInt(telegramID, 10))
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func ProblemNumbers(round Round, participant Participant) []string {
	problemIDs := round.ProblemDistribution[participant.ID]
	result := []string{}
	for i := 0; i < len(problemIDs); i++ {
		result = append(result, strconv.Itoa(i+1))
	}
	return result
}
