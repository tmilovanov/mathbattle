package mathbattle

import "time"

type Stat struct {
	//TODO: Add VisitorsToday
	//TODO: Add Current online
	ParticipantsTotal int `json:"participants_total"`
	ParticipantsToday int `json:"participants_today"`

	RoundStage       RoundStage    `json:"round_stage"`
	TimeToSolveLeft  time.Duration `json:"time_to_solve_left"`
	TimeToReviewLeft time.Duration `json:"time_to_review_left"`
	SolutionsTotal   int           `json:"solutions_total"`
	ReviewsTotal     int           `json:"reviews_total"`
	// TODO: Add SolutionsToday
	// TODO: Add ReviewsToday
}

type StatService interface {
	Stat() (Stat, error)
}
