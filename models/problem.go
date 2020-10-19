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

func ValidateProblemNumber(userInput string, problemIDs []string) (int, bool) {
	problemNumber, err := strconv.Atoi(userInput)
	if err != nil {
		return -1, false
	}
	problemNumber = problemNumber - 1
	if problemNumber < 0 || problemNumber >= len(problemIDs) {
		return -1, false
	}

	return problemNumber, true
}
