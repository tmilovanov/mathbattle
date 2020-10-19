package problemdistributor

import (
	"strconv"
	"testing"

	mathbattle "mathbattle/models"
)

type case1 struct {
	grade             int
	participantsCount int
	problemsCount     int
}

func (c *case1) Participants() []mathbattle.Participant {
	result := []mathbattle.Participant{}
	for i := 0; i < c.participantsCount; i++ {
		result = append(result, mathbattle.Participant{
			ID:    strconv.Itoa(i),
			Grade: c.grade,
		})
	}
	return result
}

func (c *case1) Problems() []mathbattle.Problem {
	result := []mathbattle.Problem{}
	for i := 0; i < c.problemsCount; i++ {
		result = append(result, mathbattle.Problem{
			ID:       strconv.Itoa(i),
			MinGrade: c.grade,
			MaxGrade: c.grade,
		})
	}
	return result
}

func (c *case1) Rounds() []mathbattle.Round {
	return []mathbattle.Round{}
}

func TestGet(t *testing.T) {
	d := RandomDistributor{}

	c := case1{grade: 7, participantsCount: 10, problemsCount: 1}
	res, _ := d.GetForAll(c.Participants(), c.Problems(), c.Rounds(), 1)
	t.Logf("%v", res)
	for _, problemIDs := range res {
		if len(problemIDs) != 1 {
			t.Fail()
		}

		if problemIDs[0] != "0" {
			t.Fail()
		}
	}

	c = case1{grade: 7, participantsCount: 10, problemsCount: 2}
	res, _ = d.GetForAll(c.Participants(), c.Problems(), c.Rounds(), 1)
	problemIDcount := make(map[string]int)
	t.Logf("%v", res)
	for _, problemIDs := range res {
		if len(problemIDs) != 1 {
			t.Fail()
		}

		problemIDcount[problemIDs[0]] += 1
	}
	t.Logf("%v", problemIDcount)

	c = case1{grade: 7, participantsCount: 10, problemsCount: 5}
	res, _ = d.GetForAll(c.Participants(), c.Problems(), c.Rounds(), 1)
	problemIDcount = make(map[string]int)
	t.Logf("%v", res)
	for _, problemIDs := range res {
		if len(problemIDs) != 1 {
			t.Fail()
		}

		problemIDcount[problemIDs[0]] += 1
	}
	t.Logf("%v", problemIDcount)
}
