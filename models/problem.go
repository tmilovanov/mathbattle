package models

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
