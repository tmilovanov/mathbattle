package sqlite_test

import (
	"testing"

	"mathbattle/repository/sqlite"

	"github.com/stretchr/testify/suite"
)

type roundTs struct {
	suite.Suite

	rep sqlite.RoundRepository
}

func (s *roundTs) SetupTest() {
	var err error

	s.rep, err = sqlite.NewRoundRepositoryTemp("mathbattle_test.sqlite")
	s.Require().Nil(err)
}

func (s *roundTs) TestStoreAndGetOne() {
	//testRound := mathbattle.Round{
	//SolveStartDate:  time.Now(),
	//SolveEndDate:    time.Now().AddDate(0, 0, 2),
	//ReviewStartDate: time.Now().AddDate(0, 0, 3),
	//ReviewEndDate:   time.Now().AddDate(0, 0, 4),
	//ProblemDistribution: map[string][]string{
	//"1": {"problem1", "problem2", "problem3"},
	//"2": {"problem1", "problem2", "problem3"},
	//"3": {"problem1", "problem2", "problem3"},
	//"4": {"problem1", "problem2", "problem3"},
	//},
	//ReviewDistribution: mathbattle.ReviewDistribution{
	//BetweenParticipants: map[string][]string{
	//"1": {"s2", "s3"},
	//"2": {"s1", "s3"},
	//"3": {"s1", "s2"},
	//},
	//ToOrganizers: []mathbattle.Solution{},
	//},
	//}
	//round, err := s.rep.Store(testRound)
	//s.Require().Nil(err)
	//s.Require().NotEqual("", round.ID)
	//testRound.ID = round.ID
	//s.Require().Equal(testRound, round)

	//round, err = s.rep.Get(round.ID)
	//s.Require().Nil(err)
	//s.Require().Equal(testRound, round)
}

func (s *roundTs) TestGetAll() {
}

func TestRoundRepository(t *testing.T) {
	suite.Run(t, &roundTs{})
}
