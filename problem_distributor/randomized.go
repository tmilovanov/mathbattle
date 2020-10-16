package problemdistributor

import (
	"fmt"
	"sort"

	mathbattle "mathbattle/models"
)

type RandomDistributor struct{}

func isFound(problems []string, problemID string) bool {
	for i := 0; i < len(problems); i++ {
		if problems[i] == problemID {
			return true
		}
	}
	return false
}

func isProblemAlreadyUsed(participant mathbattle.Participant, problem mathbattle.Problem, pastRounds []mathbattle.Round) bool {
	for _, round := range pastRounds {
		problemIDs, isExist := round.ProblemDistribution[participant.ID]
		if !isExist { // User didn't participated in this round
			continue
		}

		if isFound(problemIDs, problem.ID) {
			return true
		}
	}

	return false
}

func getSuitableProblems(participant mathbattle.Participant, problems []mathbattle.Problem,
	pastRounds []mathbattle.Round) []mathbattle.Problem {

	result := []mathbattle.Problem{}
	for _, problem := range problems {
		if !mathbattle.IsProblemSuitableForParticipant(&problem, &participant) ||
			isProblemAlreadyUsed(participant, problem, pastRounds) {
			continue
		}

		result = append(result, problem)
	}
	return result
}

type usageCounter struct {
	usage map[string]int
}

type problemUsage struct {
	id    string
	count int
}

func newUsageCounter() usageCounter {
	return usageCounter{
		usage: make(map[string]int),
	}
}

func (c *usageCounter) getSortedUsage(problems []string) []problemUsage {
	result := []problemUsage{}
	for _, id := range problems {
		_, ok := c.usage[id]
		if !ok {
			c.usage[id] = 0
		}

		result = append(result, problemUsage{
			id:    id,
			count: c.usage[id],
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].count < result[j].count
	})

	return result
}

func (c *usageCounter) getLessUsed(problems []string, count int) []string {
	usage := c.getSortedUsage(problems)
	result := []string{}
	for i := 0; i < count; i++ {
		c.usage[usage[i].id] += 1
		result = append(result, usage[i].id)
	}
	return result
}

func (d *RandomDistributor) Get(participants []mathbattle.Participant, problems []mathbattle.Problem,
	rounds []mathbattle.Round, count int) (mathbattle.RoundDistribution, error) {

	var result mathbattle.RoundDistribution = make(map[string][]string)
	counter := newUsageCounter()

	for _, participant := range participants {
		suitableProblems := getSuitableProblems(participant, problems, rounds)
		if len(suitableProblems) < count {
			return result, fmt.Errorf("No suitable problems for this participant %v", participant)
		}

		result[participant.ID] = counter.getLessUsed(mathbattle.GetProblemIDs(suitableProblems), count)
	}

	return result, nil
}
