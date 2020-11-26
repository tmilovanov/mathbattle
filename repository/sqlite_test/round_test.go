package sqlite_test

import (
	"testing"
	"time"

	mathbattle "mathbattle/models"
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

func (s *roundTs) TestSetGetUpdateDelete() {
	testRound := mathbattle.Round{
		ProblemDistribution: map[string][]string{
			"1": {"problem1", "problem2", "problem3"},
			"2": {"problem1", "problem2", "problem3"},
			"3": {"problem1", "problem2", "problem3"},
			"4": {"problem1", "problem2", "problem3"},
		},
		ReviewDistribution: mathbattle.ReviewDistribution{
			BetweenParticipants: map[string][]string{
				"1": {"s2", "s3"},
				"2": {"s1", "s3"},
				"3": {"s1", "s2"},
			},
			ToOrganizers: []string{"s4", "s5", "s6", "s7"},
		},
	}
	testRound.SetSolveStartDate(time.Now())
	testRound.SetSolveEndDate(time.Now().AddDate(0, 0, 2))
	testRound.SetReviewStartDate(time.Now().AddDate(0, 0, 3))
	testRound.SetReviewEndDate(time.Now().AddDate(0, 0, 4))

	round, err := s.rep.Store(testRound)
	s.Require().Nil(err)
	s.Require().NotEqual("", round.ID)
	testRound.ID = round.ID
	s.Require().Equal(testRound, round)

	round, err = s.rep.Get(round.ID)
	s.Require().Nil(err)
	s.Require().Equal(testRound, round)

	round.SetSolveStartDate(time.Now())
	round.SetSolveEndDate(time.Now())
	round.SetReviewStartDate(time.Now())
	round.SetReviewEndDate(time.Now())
	round.ProblemDistribution["5"] = []string{"problem1", "problem2", "problem3"}
	round.ReviewDistribution.BetweenParticipants["4"] = []string{"s5", "s6"}
	round.ReviewDistribution.ToOrganizers = append(round.ReviewDistribution.ToOrganizers, "s8", "s9", "s10")
	s.Require().Nil(s.rep.Update(round))

	updatedRound, err := s.rep.Get(round.ID)
	s.Require().Nil(err)
	s.Require().Equal(round, updatedRound)

	err = s.rep.Delete(round.ID)
	s.Require().Nil(err)

	_, err = s.rep.Get(round.ID)
	s.Require().Equal(err, mathbattle.ErrNotFound)
}

func TestRoundRepository(t *testing.T) {
	suite.Run(t, &roundTs{})
}
