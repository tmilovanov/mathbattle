package mathbattle

// Review is review on solution that one participant sends to another
type Review struct {
	ID         string `json:"id"`
	ReviewerID string `json:"reviewer_id"`
	SolutionID string `json:"solution_id"`
	Content    string `json:"content"`
}

type ReviewRepository interface {
	Store(review Review) (Review, error) // Return newly created Review with filled in ID
	Get(ID string) (Review, error)
	FindMany(reviewerID, solutionID string) ([]Review, error)
	Update(review Review) error
	Delete(ID string) error
}

type ReviewFindDescriptor struct {
	ReviewerID string `json:"reviewer_id"`
	SolutionID string `json:"solution_id"`
}

type ReviewService interface {
	Store(review Review) (Review, error)
	FindMany(descriptor ReviewFindDescriptor) ([]Review, error)
	Delete(ID string) error
	RevewStageDescriptors(participantID string) ([]SolutionDescriptor, error)
}
