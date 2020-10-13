package mock

import (
	"strconv"

	mathbattle "mathbattle/models"
)

type ProblemRepository struct {
	impl []mathbattle.Problem
}

func NewProblemRepository() ProblemRepository {
	return ProblemRepository{}
}

func (r *ProblemRepository) Store(problem mathbattle.Problem) (mathbattle.Problem, error) {
	result := problem
	result.ID = strconv.Itoa(len(r.impl))
	r.impl = append(r.impl, problem)
	return result, nil
}

func (r *ProblemRepository) GetByID(ID string) (mathbattle.Problem, error) {
	i, _ := strconv.Atoi(ID)
	return r.impl[i], nil
}

func (r *ProblemRepository) GetAll() ([]mathbattle.Problem, error) {
	return r.impl, nil
}
