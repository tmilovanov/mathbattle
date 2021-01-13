package mathbattle

type Solution struct {
	ID            string  `json:"id"`
	ParticipantID string  `json:"participant_id"`
	ProblemID     string  `json:"problem_id"`
	RoundID       string  `json:"round_id"`
	JuriComment   string  `json:"juri_comment"`
	Mark          Mark    `json:"mark"`
	Parts         []Image `json:"parts"`
}

type SolutionRepository interface {
	Store(solution Solution) (Solution, error) // Return newly created Solution with filled in ID
	Get(ID string) (Solution, error)
	Find(roundID string, participantID string, problemID string) (Solution, error)
	FindMany(roundID string, participantID string, problemID string) ([]Solution, error) //Leave IDs empty if it's not important
	FindOrCreate(roundID string, participantID string, problemID string) (Solution, error)
	AppendPart(ID string, part Image) error
	Update(solution Solution) error
	Delete(ID string) error
}

type FindDescriptor struct {
	RoundID       string `json:"round_id"`
	ParticipantID string `json:"participant_id"`
	ProblemID     string `json:"problem_id"`
}

type SolutionService interface {
	Create(solution Solution) (Solution, error)
	Get(ID string) (Solution, error)
	Find(descriptor FindDescriptor) ([]Solution, error)
	Update(solution Solution) error
	Delete(ID string) error
	AppendPart(ID string, part Image) error
	GetProblemDescriptors(participantID string) ([]ProblemDescriptor, error)
}

func SplitInGroupsByProblem(solutions []Solution) map[string][]Solution {
	result := make(map[string][]Solution)
	for _, s := range solutions {
		result[s.ProblemID] = append(result[s.ProblemID], s)
	}
	return result
}

func RoundSolutionsToString(solutions []Solution) string {
	result := ""

	for _, group := range SplitInGroupsByProblem(solutions) {
		curGroup := group[0].ProblemID + ": "
		for i := 0; i < len(group)-1; i++ {
			curGroup += group[i].ParticipantID + ", "
		}
		curGroup += group[len(group)-1].ParticipantID

		result += curGroup
	}

	return result
}
