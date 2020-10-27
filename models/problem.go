package models

import "strconv"

type Problem struct {
	ID        string
	MinGrade  int
	MaxGrade  int
	Sha256sum string
	Extension string
	Content   []byte
}

type ProblemRepository interface {
	Store(problem Problem) (Problem, error)
	GetByID(ID string) (Problem, error)
	GetAll() ([]Problem, error)
}

func GetProblemIDs(problems []Problem) []string {
	result := []string{}
	for _, problem := range problems {
		result = append(result, problem.ID)
	}
	return result
}

func ValidateIndex(userInput string, strings []string) (int, bool) {
	index, err := strconv.Atoi(userInput)
	if err != nil {
		return -1, false
	}

	index = index - 1
	if index < 0 || index >= len(strings) {
		return -1, false
	}

	return index, true
}
