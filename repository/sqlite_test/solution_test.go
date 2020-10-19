package sqlite_test

import (
	"errors"
	"os"
	"testing"

	mathbattle "mathbattle/models"
	"mathbattle/repository/sqlite"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite

	dbPath       string
	solutionPath string
	rep          sqlite.SolutionRepository
}

func (s *testSuite) SetupSuite() {
	s.dbPath = "mathbattle_test.sqlite"
	s.solutionPath = "test_solution_store"
	os.Remove(s.dbPath)
	os.RemoveAll(s.solutionPath)
	s.Require().Nil(os.Mkdir(s.solutionPath, 0777))

	rep, err := sqlite.NewSolutionRepository(s.dbPath, s.solutionPath)
	s.Require().Nil(err)
	s.Require().Nil(rep.CreateFirstTime())

	s.rep = rep
}

func (s *testSuite) TearDownSuite() {
	err := os.Remove(s.dbPath)
	s.Require().True(errors.Is(err, os.ErrNotExist) || err == nil)
	s.Require().Nil(os.RemoveAll(s.solutionPath))
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
			mathbattle.Image{Extension: ".jpg", Content: []byte("123456")},
			mathbattle.Image{Extension: ".png", Content: []byte("654321")},
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
