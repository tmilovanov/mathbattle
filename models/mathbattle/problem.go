package mathbattle

type Problem struct {
	ID        string `json:"id"`
	MinGrade  int    `json:"min_grade"`
	MaxGrade  int    `json:"max_grade"`
	Sha256sum string `json:"sha256sum"`
	Extension string `json:"extension"`
	Content   []byte `json:"content"`
}

type ProblemRepository interface {
	Store(problem Problem) (Problem, error)
	GetByID(ID string) (Problem, error)
	GetAll() ([]Problem, error)
}

type ProblemService interface {
	GetByID(ID string) (Problem, error)
}

func GetProblemIDs(problems []Problem) []string {
	result := []string{}
	for _, problem := range problems {
		result = append(result, problem.ID)
	}
	return result
}

func IsProblemSuitableForParticipant(problem *Problem, participant *Participant) bool {
	if participant.Grade >= problem.MinGrade && participant.Grade <= problem.MaxGrade {
		return true
	}
	return false
}
