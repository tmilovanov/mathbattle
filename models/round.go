package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Round struct {
	ID                  string
	SolveStartDate      time.Time // После этого времени участники могут сдавать решения
	SolveEndDate        time.Time // До этого времени участники могут сдавать решения
	ReviewStartDate     time.Time // После этого времени участники получают решения от других участников на ревью
	ReviewEndDate       time.Time // До этого времени участники могут сдавать ревью на решения других участников
	ProblemDistribution RoundDistribution
	ReviewDistribution  ReviewDistribution
}

type RoundRepository interface {
	Store(round Round) (Round, error)
	Get(ID string) (Round, error)
	GetRunning() (Round, error)
	GetSolveRunning() (Round, error)
	GetReviewPending() (Round, error)
	GetReviewRunning() (Round, error)
	GetAll() ([]Round, error)
	Update(round Round) error
	Delete(roundID string) error
}

func NewRound(solveDuration time.Duration) Round {
	solveStartTime := time.Now()
	return Round{
		ID:                  "", // ID is filled by database layer
		SolveStartDate:      solveStartTime,
		SolveEndDate:        solveStartTime.Add(solveDuration),
		ProblemDistribution: make(map[string][]string),
	}
}

// RoundDistribution is a mapping from participant ID to list of problem IDs
type RoundDistribution map[string][]string

type ReviewDistribution struct {
	BetweenParticipants map[string][]string // mapping from participantID to list of solution IDs that he got
	ToOrganizers        []Solution
}

func (d *ReviewDistribution) ToString() string {
	result := ""
	for solutionID, participantIDs := range d.BetweenParticipants {
		curSDistribution := fmt.Sprintf("%s -> ", solutionID)
		curSDistribution += strings.Join(participantIDs, ",")
		result += curSDistribution + "\n"
	}
	return result
}

// Получить порядковые номера задач, которые были разосланы участнику в этом раунде
func ProblemNumbers(round Round, participant Participant) []string {
	problemIDs := round.ProblemDistribution[participant.ID]
	result := []string{}
	for i := 0; i < len(problemIDs); i++ {
		result = append(result, strconv.Itoa(i+1))
	}
	return result
}

// Получить порядковые номера решений, которые были посланы участнику на ревью
func SolutionNumbers(round Round, participant Participant) []string {
	solutionIDs := round.ReviewDistribution.BetweenParticipants[participant.ID]

	result := []string{}
	for i := 0; i < len(solutionIDs); i++ {
		result = append(result, strconv.Itoa(i+1))
	}

	return result
}

// Remaps map[Key][]Value -> map[Value][]Key
func Remap(input map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for key, values := range input {
		for _, val := range values {
			result[val] = append(result[val], key)
		}
	}

	return result
}
