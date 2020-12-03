package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RoundStage string

const (
	StageNotStarted    RoundStage = "StageNotStarted"
	StageFinished      RoundStage = "StageFinished"
	StageSolve         RoundStage = "StageSolve"
	StageReviewPending RoundStage = "StageReviewPending"
	StageReview        RoundStage = "StageReview"
)

type Round struct {
	ID                  string
	solveStartDate      time.Time // После этого времени участники могут сдавать решения
	solveEndDate        time.Time // До этого времени участники могут сдавать решения
	reviewStartDate     time.Time // После этого времени участники получают решения от других участников на ревью
	reviewEndDate       time.Time // До этого времени участники могут сдавать ревью на решения других участников
	ProblemDistribution RoundDistribution
	ReviewDistribution  ReviewDistribution
}

func (r *Round) SetSolveStartDate(datetime time.Time) {
	r.solveStartDate = datetime.Round(0).UTC()
}

func (r *Round) GetSolveStartDate() time.Time {
	return r.solveStartDate
}

func (r *Round) SetSolveEndDate(datetime time.Time) {
	r.solveEndDate = datetime.Round(0).UTC()
}

func (r *Round) GetSolveEndDate() time.Time {
	return r.solveEndDate
}

func (r *Round) SetReviewStartDate(datetime time.Time) {
	r.reviewStartDate = datetime.Round(0).UTC()
}

func (r *Round) GetReviewStartDate() time.Time {
	return r.reviewStartDate
}

func (r *Round) SetReviewEndDate(datetime time.Time) {
	r.reviewEndDate = datetime.Round(0).UTC()
}

func (r *Round) GetReviewEndDate() time.Time {
	return r.reviewEndDate
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
	result := Round{}
	result.SetSolveStartDate(time.Now())
	result.SetSolveEndDate(result.GetSolveStartDate().Add(solveDuration))
	result.ProblemDistribution = make(map[string][]string)
	result.ReviewDistribution.BetweenParticipants = make(map[string][]string)

	return result
}

// RoundDistribution is a mapping from participant ID to list of problem IDs
type RoundDistribution map[string][]string

type ReviewDistribution struct {
	BetweenParticipants map[string][]string // mapping from participantID to list of solution IDs that he got
	ToOrganizers        []string
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

func GetRoundStage(round Round) RoundStage {
	if round.GetSolveStartDate().IsZero() || round.GetSolveStartDate().After(time.Now()) {
		return StageNotStarted
	}

	if round.GetSolveEndDate().IsZero() || round.GetSolveEndDate().After(time.Now()) {
		return StageSolve
	}

	if round.GetReviewStartDate().IsZero() || round.GetReviewStartDate().After(time.Now()) {
		return StageReviewPending
	}

	if round.GetReviewEndDate().IsZero() || round.GetReviewEndDate().After(time.Now()) {
		return StageReview
	}

	return StageFinished
}

func ParseStageEndDate(endDateTime string) (time.Time, error) {
	endDateTime = strings.Trim(endDateTime, " \t\n")
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, err
	}
	if len(endDateTime) == len("DD.MM.YYYY") {
		t, err := time.Parse("02.01.2006", endDateTime)
		if err != nil {
			return time.Time{}, ErrWrongUserInput
		}
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, moscowLocation)
		t = t.AddDate(0, 0, 1)
		return t, nil
	} else if len(endDateTime) == len("DD.MM.YYYY HH:MM") {
		t, err := time.Parse("02.01.2006 15:04", endDateTime)
		if err != nil {
			return time.Time{}, ErrWrongUserInput
		}
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, moscowLocation)
		return t, nil
	}

	return time.Time{}, ErrWrongUserInput
}
