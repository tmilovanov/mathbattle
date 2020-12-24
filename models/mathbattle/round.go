package mathbattle

import (
	"errors"
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
	ID string `json:"id"`

	// После этого времени участники могут сдавать решения
	SolveStartDate time.Time `json:"solve_start_date"`

	// До этого времени участники могут сдавать решения
	SolveEndDate time.Time `json:"solve_end_date"`

	// После этого времени участники получают решения от других участников на ревью
	ReviewStartDate time.Time `json:"review_start_date"`

	// До этого времени участники могут сдавать ревью на решения других участников
	ReviewEndDate time.Time `json:"review_end_date"`

	ProblemDistribution RoundDistribution  `json:"problem_distribution"`
	ReviewDistribution  ReviewDistribution `json:"solution_distribution"`
}

func (r *Round) IsActive() bool {
	stage := GetRoundStage(*r)
	if stage == StageFinished || stage == StageNotStarted {
		return false
	}

	return true
}

func (r *Round) SetSolveStartDate(datetime time.Time) {
	r.SolveStartDate = datetime.Round(0).UTC()
}

func (r *Round) GetSolveStartDate() time.Time {
	return r.SolveStartDate
}

func (r *Round) SetSolveEndDate(datetime time.Time) {
	r.SolveEndDate = datetime.Round(0).UTC()
}

func (r *Round) GetSolveEndDate() time.Time {
	return r.SolveEndDate
}

func (r *Round) GetSolveStageDuration() time.Duration {
	return r.GetSolveEndDate().Sub(r.GetSolveStartDate())
}

func (r *Round) GetSolveEndDateMsk() (time.Time, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, err
	}
	return r.SolveEndDate.In(location), nil
}

func (r *Round) SetReviewStartDate(datetime time.Time) {
	r.ReviewStartDate = datetime.Round(0).UTC()
}

func (r *Round) GetReviewStartDate() time.Time {
	return r.ReviewStartDate
}

func (r *Round) SetReviewEndDate(datetime time.Time) {
	r.ReviewEndDate = datetime.Round(0).UTC()
}

func (r *Round) GetReviewEndDate() time.Time {
	return r.ReviewEndDate
}

func (r *Round) GetReviewEndDateMsk() (time.Time, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, err
	}
	return r.ReviewEndDate.In(location), nil
}

func (r *Round) GetReviewStageDuration() time.Duration {
	return r.GetReviewEndDate().Sub(r.GetReviewStartDate())
}

type RoundRepository interface {
	Store(round Round) (Round, error)
	Get(ID string) (Round, error)
	GetRunning() (Round, error)
	GetSolveRunning() (Round, error)
	GetReviewPending() (Round, error)
	GetReviewRunning() (Round, error)
	GetAll() ([]Round, error)
	GetLast() (Round, error)
	Update(round Round) error
	Delete(roundID string) error
}

type StartOrder struct {
	StageEnd string `json:"stage_end"`
}

type RoundService interface {
	StartNew(startOrder StartOrder) (Round, error)
	StartReviewStage(startOrder StartOrder) (Round, error)
	ReviewStageDistributionDesc() (ReviewDistributionDesc, error)
	GetAll() ([]Round, error)
	GetByID(ID string) (Round, error)
	GetRunning() (Round, error)
	GetReviewPending() (Round, error)
	GetReviewRunning() (Round, error)
	GetLast() (Round, error)
	GetProblemDescriptors(participantID string) ([]ProblemDescriptor, error)
}

func NewRoundFromEnd(solveStageEnd time.Time) Round {
	result := Round{}
	result.SetSolveStartDate(time.Now())
	result.SetSolveEndDate(solveStageEnd)
	result.ProblemDistribution = make(map[string][]ProblemDescriptor)
	result.ReviewDistribution.BetweenParticipants = make(map[string][]string)

	return result
}

func NewRoundFromDuration(solveDuration time.Duration) Round {
	endTime := time.Now().Add(solveDuration)
	return NewRoundFromEnd(endTime)
}

type ReviewDistributionDesc struct {
	Desc string `json:"desc"`
}

type ProblemDescriptor struct {
	// Caption используется вместо названия задачи. Нужен только для того чтобы конкретный участник мог
	// как-то назвать задачу. Например "задача А" или "задача 1"
	// Уникален для каждой задачи при фиксированном Round и Participant
	Caption string `json:"caption"`
	// ID из базы
	ProblemID string `json:"problem_id"`
}

// RoundDistribution is a mapping from participant ID to list of problems that participant get to solve
type RoundDistribution map[string][]ProblemDescriptor

func (pd RoundDistribution) FindDescriptor(participantID string, problemID string) (ProblemDescriptor, error) {
	descriptors, isExist := pd[participantID]
	if !isExist {
		return ProblemDescriptor{}, ErrNotFound
	}

	for i := 0; i < len(descriptors); i++ {
		if descriptors[i].ProblemID == problemID {
			return descriptors[i], nil
		}
	}

	return ProblemDescriptor{}, ErrNotFound
}

type ReviewDistribution struct {
	// Mapping from participantID to list of solution IDs that he got
	BetweenParticipants map[string][]string `json:"between_participants"`
	ToOrganizers        []string            `json:"to_organizers"`
}

func ProblemIDsFromSolutionIDs(solutions SolutionRepository, solutionIDs []string) ([]string, error) {
	result := []string{}
	for _, ID := range solutionIDs {
		solution, err := solutions.Get(ID)
		if err != nil {
			return result, err
		}

		result = append(result, solution.ProblemID)
	}

	return result, nil
}

type SolutionDescriptor struct {
	ProblemCaption string
	SolutionNumber int
	SolutionID     string
}

func SolutionDescriptorsForParticipant(problemIDs []string, solutionIDs []string,
	participantID string, round Round) ([]SolutionDescriptor, error) {

	result := []SolutionDescriptor{}

	solutionNumbers := make(map[string]int)
	if len(problemIDs) != len(solutionIDs) {
		return result, errors.New("Expect the same count for problemIDs and solutionIDs")
	}

	for i := 0; i < len(problemIDs); i++ {
		desc, err := round.ProblemDistribution.FindDescriptor(participantID, problemIDs[i])
		if err != nil {
			return result, err
		}

		solutionNumbers[desc.Caption]++

		result = append(result, SolutionDescriptor{
			ProblemCaption: desc.Caption,
			SolutionNumber: solutionNumbers[desc.Caption],
			SolutionID:     solutionIDs[i],
		})

	}

	return result, nil
}

func SolutionDescriptorsFromSolutionIDs(solutions SolutionRepository,
	participantID string, round Round) ([]SolutionDescriptor, error) {

	solutionIDs := round.ReviewDistribution.BetweenParticipants[participantID]
	problemIDs, err := ProblemIDsFromSolutionIDs(solutions, solutionIDs)
	if err != nil {
		return []SolutionDescriptor{}, err
	}

	return SolutionDescriptorsForParticipant(problemIDs, solutionIDs, participantID, round)
}

func FindSolutionIDbyDescriptor(item SolutionDescriptor, descriptors []SolutionDescriptor) (string, bool) {
	for _, descriptor := range descriptors {
		if descriptor.ProblemCaption == item.ProblemCaption &&
			descriptor.SolutionNumber == item.SolutionNumber {
			return descriptor.SolutionID, true
		}
	}

	return "", false
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

// Получить обозначения задач, которые были разосланы участнику в этом раунде
func ProblemsCaptions(descriptors []ProblemDescriptor) []string {
	result := []string{}
	for _, descriptor := range descriptors {
		result = append(result, descriptor.Caption)
	}
	return result
}

// Получить обозначения задач, которые были разосланы участнику в этом раунде
func ProblemsCaptionsOld(round Round, participant Participant) []string {
	result := []string{}
	for _, descriptor := range round.ProblemDistribution[participant.ID] {
		result = append(result, descriptor.Caption)
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

func ValidateCaptions(userInput string, descriptors []ProblemDescriptor) (int, bool) {
	userInput = strings.Trim(userInput, "\t\r\n ")
	for i, desc := range descriptors {
		if desc.Caption == userInput {
			return i, true
		}
	}

	return -1, false
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
