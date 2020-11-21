package models

// Review is review on solution that one participant sends to another
type Review struct {
	ID         string
	ReviewerID string
	SolutionID string
	Content    string
}

type ReviewRepository interface {
	Store(review Review) (Review, error) // Return newly created Review with filled in ID
	Get(ID string) (Review, error)
	FindMany(reviewerID, solutionID string) ([]Review, error)
	Update(review Review) error
	Delete(ID string) error
}
