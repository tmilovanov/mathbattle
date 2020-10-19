package models

type Solution struct {
	ID            string
	ParticipantID string
	ProblemID     string
	RoundID       string
	Parts         []Image
}

type SolutionRepository interface {
	Store(solution Solution) (Solution, error) // Return newly created Solution with filled in ID
	Get(ID string) (Solution, error)
	Find(roundID string, participantID string, problemID string) (Solution, error)
	FindMany(roundID string, participantID string, problemID string) ([]Solution, error) //Leave IDs empty if it's not important
	FindOrCreate(roundID string, participantID string, problemID string) (Solution, error)
	AppendPart(ID string, part Image) error
	Delete(ID string) error
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
