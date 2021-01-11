package repositorytest

import (
	"log"
	"testing"

	"mathbattle/infrastructure"
	"mathbattle/mocks"
	"mathbattle/models/mathbattle"

	"github.com/stretchr/testify/suite"
)

type problemTs struct {
	suite.Suite

	rep mathbattle.ProblemRepository
}

func (s *problemTs) SetupTest() {
	container := infrastructure.NewTestContainer()
	s.rep = container.ProblemRepository()
}

func (s *problemTs) TestStore() {
	log.Printf("TestStore")
	for _, problem := range mocks.GenProblems(10, 1, 11) {
		p, err := s.rep.Store(problem)
		s.Require().Nil(err)

		problem.ID = p.ID
		problem.Sha256sum = p.Sha256sum
		s.Require().Equal(problem, p)
	}
}

func (s *problemTs) TestGetAll() {
	log.Printf("TestGetAll")
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
