package repositorytest

import (
	"testing"

	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	"github.com/stretchr/testify/suite"
)

type solutionTs struct {
	suite.Suite

	rep mathbattle.SolutionRepository
}

func (s *solutionTs) SetupSuite() {
	container := infrastructure.NewTestContainer()
	s.rep = container.SolutionRepository()
}

func (s *solutionTs) TestCreateNewEmpty() {
	newEmptySolution := mathbattle.Solution{RoundID: "1", ParticipantID: "1", ProblemID: "1"}
	solution, err := s.rep.Store(newEmptySolution)
	newEmptySolution.ID = solution.ID
	s.Require().Nil(err)
	s.Require().Equal(newEmptySolution, solution)
	newEmptySolution = solution

	solution, err = s.rep.Find("1", "1", "1")
	s.Require().Nil(err)
	s.Require().Equal(newEmptySolution, solution)

	err = s.rep.Delete(solution.ID)
	s.Require().Nil(err)

	solution, err = s.rep.Find("1", "1", "1")
	s.Require().Equal(err, mathbattle.ErrNotFound)
}

func (s *solutionTs) TestFindMany() {
	allSolutions := []mathbattle.Solution{
		{
			RoundID:       "1",
			ParticipantID: "1",
			ProblemID:     "1",
			Parts: []mathbattle.Image{
				{Extension: ".jpg", Content: []byte("1111111")},
				{Extension: ".png", Content: []byte("2222222")},
			},
		},
		{
			RoundID:       "1",
			ParticipantID: "2",
			ProblemID:     "2",
			Parts: []mathbattle.Image{
				{Extension: ".jpg", Content: []byte("33333333")},
				{Extension: ".png", Content: []byte("44444444")},
			},
		},
	}

	for i, solution := range allSolutions {
		sl, err := s.rep.Store(solution)
		s.Require().Nil(err)
		allSolutions[i].ID = sl.ID
		s.Require().Equal(allSolutions[i], sl)
	}

	solutions, err := s.rep.FindMany("1", "", "")
	s.Require().Nil(err)
	s.Require().Equal(2, len(solutions))
	s.Require().Equal(allSolutions, solutions)

	solutions, err = s.rep.FindMany("1", "", "1")
	s.Require().Nil(err)
	s.Require().Equal(1, len(solutions))
	s.Require().Equal(allSolutions[0], solutions[0])

	solutions, err = s.rep.FindMany("1", "", "2")
	s.Require().Nil(err)
	s.Require().Equal(1, len(solutions))
	s.Require().Equal(allSolutions[1], solutions[0])
}

func (s *solutionTs) TestCreateNewNotEmptyAndAppend() {
	newSolution := mathbattle.Solution{
		RoundID:       "1",
		ParticipantID: "1",
		ProblemID:     "1",
		Parts: []mathbattle.Image{
			{Extension: ".jpg", Content: []byte("123456")},
			{Extension: ".png", Content: []byte("654321")},
		},
		JuriComment: "good solution",
		Mark:        10,
	}
	solution, err := s.rep.Store(newSolution)
	newSolution.ID = solution.ID
	s.Require().Nil(err)
	s.Require().Equal(newSolution, solution)
	newSolution = solution

	solution, err = s.rep.Find("1", "1", "1")
	s.Require().Nil(err)
	s.Require().Equal(newSolution, solution)

	newPart := mathbattle.Image{Extension: ".jpg", Content: []byte("55555")}
	newSolution.Parts = append(newSolution.Parts, newPart)

	err = s.rep.AppendPart(solution.ID, newPart)
	s.Require().Nil(err)

	solution, err = s.rep.Get(solution.ID)
	s.Require().Nil(err)
	s.Require().Equal(solution, newSolution)

	s.Require().Nil(s.rep.Delete(solution.ID))
	solution, err = s.rep.Get(solution.ID)
	s.Require().Equal(err, mathbattle.ErrNotFound)
}

func TestRepository(t *testing.T) {
	suite.Run(t, &solutionTs{})
}
