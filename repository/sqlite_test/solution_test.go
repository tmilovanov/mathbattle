package sqlite_test

import (
	"errors"
	"os"
	"testing"

	mathbattle "mathbattle/models"
	"mathbattle/repository/sqlite"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	dbPath       string
	solutionPath string
	rep          sqlite.SQLSolutionRepository
	req          *require.Assertions
}

func (s *testSuite) SetupSuite() {
	s.req = require.New(s.T())
	s.dbPath = "mathbattle_test.sqlite"
	s.solutionPath = "test_solution_store"
	os.Remove(s.dbPath)
	os.RemoveAll(s.solutionPath)
	s.req.Nil(os.Mkdir(s.solutionPath, 0777))

	rep, err := sqlite.NewSQLSolutionRepository(s.dbPath, s.solutionPath)
	s.req.Nil(err)
	s.req.Nil(rep.CreateFirstTime())

	s.rep = rep
}

func (s *testSuite) TearDownSuite() {
	err := os.Remove(s.dbPath)
	s.req.True(errors.Is(err, os.ErrNotExist) || err == nil)
	s.req.Nil(os.RemoveAll(s.solutionPath))
}

func (s *testSuite) TestCreateNewEmpty() {
	newEmptySolution := mathbattle.Solution{RoundID: "1", ParticipantID: "1", ProblemID: "1"}
	solution, err := s.rep.Store(newEmptySolution)
	newEmptySolution.ID = solution.ID
	s.req.Nil(err)
	s.req.Equal(newEmptySolution, solution)
	newEmptySolution = solution

	solution, err = s.rep.Find("1", "1", "1")
	s.req.Nil(err)
	s.req.Equal(newEmptySolution, solution)

	err = s.rep.Delete(solution.ID)
	s.req.Nil(err)

	solution, err = s.rep.Find("1", "1", "1")
	s.req.Equal(err, mathbattle.ErrNotFound)
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
	s.req.Nil(err)
	s.req.Equal(newSolution, solution)
	newSolution = solution

	solution, err = s.rep.Find("1", "1", "1")
	s.req.Nil(err)
	s.req.Equal(newSolution, solution)

	newPart := mathbattle.Image{Extension: ".jpg", Content: []byte("55555")}
	newSolution.Parts = append(newSolution.Parts, newPart)

	err = s.rep.AppendPart(solution.ID, newPart)
	s.req.Nil(err)

	solution, err = s.rep.Get(solution.ID)
	s.req.Nil(err)
	s.req.Equal(solution, newSolution)

	s.req.Nil(s.rep.Delete(solution.ID))
	solution, err = s.rep.Get(solution.ID)
	s.req.Equal(err, mathbattle.ErrNotFound)
}

func TestRepository(t *testing.T) {
	suite.Run(t, &testSuite{})
}
