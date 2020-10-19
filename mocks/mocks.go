package mocks

import (
	"fmt"
	"log"
	"time"

	"mathbattle/combinator"
	mathbattle "mathbattle/models"
)

func GenProblems(problemsCount int) []mathbattle.Problem {
	result := []mathbattle.Problem{}

	for i := 0; i < problemsCount; i++ {
		result = append(result, mathbattle.Problem{
			MinGrade:  1,
			MaxGrade:  11,
			Extension: ".jpg",
			Content:   []byte(fmt.Sprintf("%d fake problem", i)),
		})
	}

	return result
}

func GenParticipants(participantsCount int) []mathbattle.Participant {
	result := []mathbattle.Participant{}

	for i := 0; i < participantsCount; i++ {
		result = append(result, mathbattle.Participant{
			//TelegramID: fmt.Sprintf("%d fake telegram id", i),
			Name:  fmt.Sprintf("%d fake name", i),
			Grade: 5,
		})
	}

	return result
}

func Solve(solutions mathbattle.SolutionRepository, problemSolvesCount []int) {
	solutions.Store(mathbattle.Solution{
		ParticipantID: "",
		ProblemID:     "",
		RoundID:       "",
	})
}

func GenSolutionStageRound(rounds mathbattle.RoundRepository, participants mathbattle.ParticipantRepository,
	problems mathbattle.ProblemRepository, problemDistributor mathbattle.ProblemDistributor,
	participantsCount int, problemOnEach int) (mathbattle.Round, error) {

	allProblems := GenProblems(problemOnEach)
	for _, problem := range allProblems {
		_, err := problems.Store(problem)
		if err != nil {
			return mathbattle.Round{}, err
		}
	}

	allParticipants := GenParticipants(participantsCount)
	for _, participant := range allParticipants {
		_, err := participants.Store(participant)
		if err != nil {
			return mathbattle.Round{}, err
		}
	}

	distribution, err := problemDistributor.Get(allParticipants, allProblems, []mathbattle.Round{})
	if err != nil {
		return mathbattle.Round{}, err
	}

	round, err := rounds.Store(mathbattle.Round{
		SolveStartDate:      time.Now(),
		ProblemDistribution: distribution,
	})

	return round, err
}

func GenReviewPendingRound(rounds mathbattle.RoundRepository, participants mathbattle.ParticipantRepository,
	solutions mathbattle.SolutionRepository, problems mathbattle.ProblemRepository, problemDistributor mathbattle.ProblemDistributor,
	participantsCount int, problemOnEach int) (mathbattle.Round, error) {

	round, err := GenSolutionStageRound(rounds, participants, problems,
		problemDistributor, participantsCount, problemOnEach)
	if err != nil {
		return round, err
	}

	//round.SolveEndDate := time.Now().AddDate(0,0,-1)

	return mathbattle.Round{}, nil
}

func GenProblemIDs(problemCount int) []string {
	result := []string{}
	for i := 0; i < problemCount; i++ {
		id := int('A') + i
		if i >= 'Z' {
			log.Panic("problemCount is too large")
		}
		result = append(result, string(rune(id)))
	}
	return result
}

func GenAllSolutionsCombinations(problemCount, participantCount int) [][]mathbattle.Solution {
	result := [][]mathbattle.Solution{}
	for _, combination := range combinator.GetAll(problemCount, participantCount) {
		result = append(result, genOneSolutionCombination(combination))
	}
	return result
}

func genOneSolutionCombination(solutionsCount []int) []mathbattle.Solution {
	result := []mathbattle.Solution{}
	problemCount := len(solutionsCount)
	problemIDs := GenProblemIDs(problemCount)
	for i := 0; i < problemCount; i++ {
		for j := 0; j < solutionsCount[i]; j++ {
			pariticipantID := fmt.Sprintf("p%d", j)
			result = append(result, mathbattle.Solution{
				ID:            fmt.Sprintf("s_%s_%s", pariticipantID, problemIDs[i]),
				ParticipantID: pariticipantID,
				ProblemID:     problemIDs[i],
			})
		}
	}
	return result
}
