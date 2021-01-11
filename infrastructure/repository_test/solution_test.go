package repositorytest

import (
	"testing"

	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite

	rep mathbattle.SolutionRepository
}

func (s *testSuite) SetupSuite() {
	container := infrastructure.NewTestContainer()
	s.rep = container.SolutionRepository()
}

func (s *testSuite) TestCreateNewEmpty() {
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

func (s *testSuite) TestCreateNewNotEmptyAndAppend() {
	newSolution := mathbattle.Solution{
		RoundID:       "1",
		ParticipantID: "1",
		ProblemID:     "1",
		Parts: []mathbattle.Image{
			{Extension: ".jpg", Content: []byte("123456")},
			{Extension: ".png", Content: []byte("654321")},
		},
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
	suite.Run(t, &testSuite{})
}
