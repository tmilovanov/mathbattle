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
	if descriptor.ProblemID == "" {
		return s.Rep.FindMany(descriptor.ReviewerID, descriptor.SolutionID)
	} else {
		result := []mathbattle.Review{}

		allSolutions, err := s.Solutions.FindMany("", descriptor.SolutionID, descriptor.ProblemID)
		if err != nil {
			return result, err
		}

		if len(allSolutions) == 0 {
			return result, nil
		}

		for _, solution := range allSolutions {
			reviews, err := s.Rep.FindMany(descriptor.ReviewerID, solution.ID)
			if err != nil {
				return []mathbattle.Review{}, err
			}
			result = append(result, reviews...)
		}

		return result, nil
	}
}

func (s *ReviewService) Delete(ID string) error {
	return s.Rep.Delete(ID)
}

func (s *ReviewService) RevewStageDescriptors(participantID string) ([]mathbattle.SolutionDescriptor, error) {
	round, err := s.Rounds.GetLast()
	if err != nil {
		return []mathbattle.SolutionDescriptor{}, err
	}

	return mathbattle.SolutionDescriptorsFromSolutionIDs(s.Solutions, participantID, round)
}
