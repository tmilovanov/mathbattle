package sqlitetest

import (
	"testing"

	"mathbattle/infrastructure/repository/sqlite"
	"mathbattle/mocks"

	"github.com/stretchr/testify/suite"
)

type problemTs struct {
	suite.Suite

	rep *sqlite.ProblemRepository
}

func (s *problemTs) SetupTest() {
	var err error

	DeleteTempDatabase()
	s.rep, err = sqlite.NewProblemRepository(TestDbPath(), TestProblemsPath())
	s.Require().Nil(err)
}

func (s *problemTs) TestStore() {
	for _, problem := range mocks.GenProblems(10, 1, 11) {
		p, err := s.rep.Store(problem)
		s.Require().Nil(err)

		problem.ID = p.ID
		problem.Sha256sum = p.Sha256sum
		s.Require().Equal(problem, p)
	}
}

func (s *problemTs) TestGetAll() {
	var err error

	problems := mocks.GenProblems(10, 1, 11)
	for i := 0; i < len(problems); i++ {
		problems[i], err = s.rep.Store(problems[i])
		s.Require().Nil(err)
	}

	storedProblems, err := s.rep.GetAll()
	s.Require().Nil(err)
	s.Require().Equal(problems, storedProblems)
}

func TestProblemsRepository(t *testing.T) {
	suite.Run(t, &problemTs{})
}
