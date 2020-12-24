package application

import "mathbattle/models/mathbattle"

type ReviewService struct {
	Rep       mathbattle.ReviewRepository
	Rounds    mathbattle.RoundRepository
	Solutions mathbattle.SolutionRepository
}

func (s *ReviewService) Store(review mathbattle.Review) (mathbattle.Review, error) {
	return s.Rep.Store(review)
}

func (s *ReviewService) FindMany(descriptor mathbattle.ReviewFindDescriptor) ([]mathbattle.Review, error) {
	return s.Rep.FindMany(descriptor.ReviewerID, descriptor.SolutionID)
}

func (s *ReviewService) Delete(ID string) error {
	return s.Rep.Delete(ID)
}

func (s *ReviewService) RevewStageDescriptors(participantID string) ([]mathbattle.SolutionDescriptor, error) {
	round, err := s.Rounds.GetReviewRunning()
	if err != nil {
		return []mathbattle.SolutionDescriptor{}, err
	}

	return mathbattle.SolutionDescriptorsFromSolutionIDs(s.Solutions, participantID, round)
}
