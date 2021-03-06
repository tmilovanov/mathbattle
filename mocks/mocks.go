package mocks

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	problemdistributor "mathbattle/application/problem_distributor"
	solutiondistributor "mathbattle/application/solution_distributor"
	"mathbattle/libs/combinator"
	"mathbattle/libs/mstd"
	"mathbattle/models/mathbattle"
)

func GenProblems(problemsCount int, minGrade int, maxGrade int) []mathbattle.Problem {
	result := []mathbattle.Problem{}

	for i := 0; i < problemsCount; i++ {
		problemContent := []byte(fmt.Sprintf("%d fake problem", i))
		h := sha256.New()

		_, err := io.Copy(h, bytes.NewReader(problemContent))
		if err != nil {
			log.Fatal(err)
		}

		sha256sum := hex.EncodeToString(h.Sum(nil))

		result = append(result, mathbattle.Problem{
			Sha256sum: sha256sum,
			MinGrade:  minGrade,
			MaxGrade:  maxGrade,
			Extension: ".jpg",
			Content:   problemContent,
		})
	}

	return result
}

func GenParticipants(participantsCount int, grade int) []mathbattle.Participant {
	result := []mathbattle.Participant{}

	for i := 0; i < participantsCount; i++ {
		result = append(result, mathbattle.Participant{
			User: mathbattle.User{
				TelegramID:       int64(i),
				TelegramUsername: fmt.Sprintf("fake_telegram_name_%d", i),
				IsAdmin:          false,
				RegistrationTime: mstd.UTCNowRoundSeconds(),
			},
			Name:   fmt.Sprintf("%d fake name", i),
			School: fmt.Sprintf("Fake school %d", i),
			Grade:  grade,
		})
	}

	return result
}

// GenSolutionsForRound генерирует заданное количество решений для каждой задачи в раунде
func GenSolutionsForRound(roundID string, rd mathbattle.RoundDistribution, needSolutionsCount map[string]int) []mathbattle.Solution {
	curSolutionsCount := make(map[string]int)
	for problemID := range needSolutionsCount {
		curSolutionsCount[problemID] = 0
	}

	result := []mathbattle.Solution{}
	for participantID, problemDescriptors := range rd {
		for _, desc := range problemDescriptors {
			if curSolutionsCount[desc.ProblemID] < needSolutionsCount[desc.ProblemID] {
				result = append(result, mathbattle.Solution{
					RoundID:       roundID,
					ProblemID:     desc.ProblemID,
					ParticipantID: participantID,
					Parts: []mathbattle.Image{
						{
							Extension: ".jpg",
							Content:   []byte(fmt.Sprintf("s_of_%s_on_%s", participantID, desc.ProblemID)),
						},
					},
				})
				curSolutionsCount[desc.ProblemID]++
			}
		}
	}

	return result
}

func GenSolutionStageRound(rounds mathbattle.RoundRepository, participants mathbattle.ParticipantRepository,
	problems mathbattle.ProblemRepository, problemDistributor problemdistributor.SimpleDistributor,
	participantsCount int, problemOnEach int) (mathbattle.Round, error) {

	var err error

	allProblems := GenProblems(problemOnEach, 1, 11)
	for i := 0; i < len(allProblems); i++ {
		allProblems[i], err = problems.Store(allProblems[i])
		if err != nil {
			return mathbattle.Round{}, err
		}
	}

	allParticipants := GenParticipants(participantsCount, 5)
	for i := 0; i < len(allParticipants); i++ {
		allParticipants[i], err = participants.Store(allParticipants[i])
		if err != nil {
			return mathbattle.Round{}, err
		}
	}

	round := mathbattle.Round{
		ProblemDistribution: make(map[string][]mathbattle.ProblemDescriptor),
	}
	round.SetSolveStartDate(time.Now())

	for _, participant := range allParticipants {
		problems, err := problemDistributor.GetForParticipantCount(participant, problemOnEach)
		if err != nil {
			return mathbattle.Round{}, err
		}

		descriptors := []mathbattle.ProblemDescriptor{}
		for i, problem := range problems {
			descriptors = append(descriptors, mathbattle.ProblemDescriptor{
				Caption:   mstd.IndexToLetter(i),
				ProblemID: problem.ID,
			})
		}

		round.ProblemDistribution[participant.ID] = descriptors
	}

	round, err = rounds.Store(round)

	return round, err
}

func GenReviewPendingRound(rounds mathbattle.RoundRepository, participants mathbattle.ParticipantRepository,
	solutions mathbattle.SolutionRepository, problems mathbattle.ProblemRepository, problemDistributor problemdistributor.SimpleDistributor,
	participantsCount int, problemOnEach int, solutionsCount []int) (mathbattle.Round, error) {

	round, err := GenSolutionStageRound(rounds, participants, problems,
		problemDistributor, participantsCount, problemOnEach)
	if err != nil {
		return round, err
	}

	problemSolutionsCount := make(map[string]int)
	allProblems, err := problems.GetAll()
	if err != nil {
		return round, err
	}
	if len(solutionsCount) != len(allProblems) {
		return round, errors.New("Expect count of solutions to be equal count of problems")
	}
	for i := 0; i < len(solutionsCount); i++ {
		problemSolutionsCount[allProblems[i].ID] = solutionsCount[i]
	}

	allRoundSolutions := GenSolutionsForRound(round.ID, round.ProblemDistribution, problemSolutionsCount)
	for _, curSolution := range allRoundSolutions {
		_, err := solutions.Store(curSolution)
		if err != nil {
			return mathbattle.Round{}, err
		}
	}
	round.SetSolveEndDate(time.Now().AddDate(0, 0, -1))
	rounds.Update(round)

	return round, nil
}

func GenReviewStageRound(rounds mathbattle.RoundRepository, participants mathbattle.ParticipantRepository,
	solutions mathbattle.SolutionRepository, problems mathbattle.ProblemRepository,
	problemDistributor problemdistributor.SimpleDistributor, solutionsDistributor solutiondistributor.SolutionDistributor,
	participantsCount int, problemOnEach int, solutionsCount []int, reviewersCount uint) (mathbattle.Round, error) {

	round, err := GenReviewPendingRound(rounds, participants, solutions, problems, problemDistributor,
		participantsCount, problemOnEach, solutionsCount)
	if err != nil {
		return round, err
	}

	allRoundSolutions, err := solutions.FindMany(round.ID, "", "")
	if err != nil {
		return round, err
	}

	round.SetReviewStartDate(time.Now())
	round.ReviewDistribution = solutionsDistributor.Get(allRoundSolutions, reviewersCount)
	err = rounds.Update(round)
	if err != nil {
		return round, err
	}

	return round, nil
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
