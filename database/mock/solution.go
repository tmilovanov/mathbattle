package mock

import (
	mathbattle "mathbattle/models"
)

type SolutionRepository struct {
	impl map[string]mathbattle.Solution
}

func NewSolutionRepository() SolutionRepository {
	return SolutionRepository{make(map[string]mathbattle.Solution)}
}

func (r *SolutionRepository) getStrIDimpl(participantID string, roundID string, problemID string) string {
	return participantID + roundID + problemID
}

func (r *SolutionRepository) getStrID(solution mathbattle.Solution) string {
	return r.getStrIDimpl(solution.ParticipantID, solution.RoundID, solution.ProblemID)
}

func (r *SolutionRepository) Store(solution mathbattle.Solution) (mathbattle.Solution, error) {
	solution.ID = r.getStrID(solution)
	r.impl[solution.ID] = solution
	return solution, nil
}

func (r *SolutionRepository) Get(ID string) (mathbattle.Solution, error) {
	res, ok := r.impl[ID]
	if !ok {
		return res, mathbattle.ErrNotFound
	}
	return res, nil
}

func (r *SolutionRepository) Find(roundID string, participantID string, problemID string) (mathbattle.Solution, error) {
	res, ok := r.impl[r.getStrIDimpl(participantID, roundID, problemID)]
	if !ok {
		return res, mathbattle.ErrNotFound
	}
	return res, nil
}

func (r *SolutionRepository) FindOrCreate(roundID string, participantID string, problemID string) (mathbattle.Solution, error) {
	res, err := r.Find(roundID, participantID, problemID)
	if err != nil {
		return res, nil
	}

	return r.Store(mathbattle.Solution{
		RoundID:       roundID,
		ParticipantID: participantID,
		ProblemID:     problemID,
	})
}

func (r *SolutionRepository) AppendPart(ID string, item mathbattle.Image) error {
	res, ok := r.impl[ID]
	if !ok {
		return mathbattle.ErrNotFound
	}
	res.Parts = append(res.Parts, item)
	r.impl[ID] = res
	return nil
}

func (r *SolutionRepository) Delete(ID string) error {
	delete(r.impl, ID)
	return nil
}
