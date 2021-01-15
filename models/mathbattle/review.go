package mathbattle

// Review is review on solution that one participant sends to another
type Review struct {
	ID          string `json:"id"`
	ReviewerID  string `json:"reviewer_id"`
	SolutionID  string `json:"solution_id"`
	Content     string `json:"content"`
	JuriComment string `json:"juri_comment"`
	Mark        Mark   `json:"mark"`
}

type ReviewRepository interface {
	Store(review Review) (Review, error) // Return newly created Review with filled in ID
	Get(ID string) (Review, error)
	FindMany(reviewerID, solutionID string) ([]Review, error)
	Update(review Review) error
	Delete(ID string) error
}

// ReviewerID="", SolutionID="", ProblemID="" - Get all reviews in all rounds
// ReviewerID="", SolutionID="", ProblemID=ID - Get all reviews on all solutions of the problem with id=ID
// ReviewerID="", SolutionID=ID, ProblemID="" - Get all reviews on solution with ID=ID
// ReviewerID="", SolutionID=ID, ProblemID=ID - The same as above
// ReviewerID=ID, SolutionID="", ProblemID="" - Get all reviews in all rounds of this reviewer
// ReviewerID=ID, SolutionID="", ProblemID="" - Get all reviews in all rounds of this reviewer on all solutions of the problem with id=ID
// ReviewerID=ID, SolutionID=ID, ProblemID="" - Get all reviews in all rounds of this reviewer on solution with ID=ID
// ReviewerID=ID, SolutionID=ID, ProblemID=ID - The same as above
type ReviewFindDescriptor struct {
	ReviewerID string `json:"reviewer_id"`
	SolutionID string `json:"solution_id"`
	ProblemID  string `json:"problem_id"`
}

type ReviewService interface {
	Store(review Review) (Review, error)
	FindMany(descriptor ReviewFindDescriptor) ([]Review, error)
	Delete(ID string) error
	RevewStageDescriptors(participantID string) ([]SolutionDescriptor, error)
}
