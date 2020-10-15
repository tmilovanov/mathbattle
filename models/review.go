package models

// Review is review on solution that one participant sends to another
type Review struct {
	ID                    string
	ReviewerID            string
	ReviewedParticipantID string
	SolutionID            string
	Content               string
}
