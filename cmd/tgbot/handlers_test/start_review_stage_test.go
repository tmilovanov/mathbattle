package handlerstest

import (
	"testing"

	"mathbattle/cmd/tgbot/handlers"
	mreplier "mathbattle/cmd/tgbot/replier"
	"mathbattle/mocks"
	problemdistributor "mathbattle/problem_distributor"
	"mathbattle/repository/sqlite"
	solutiondistributor "mathbattle/solution_distributor"

	"github.com/stretchr/testify/suite"
)

type startReviewStageTs struct {
	suite.Suite

	handler handlers.StartReviewStage
	replier mreplier.Replier
	chatID  int64
}

func (s *startReviewStageTs) SetupTest() {
	s.chatID = 123456
	s.replier = mreplier.RussianReplier{}
	participants, err := sqlite.NewParticipantRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	solutions, err := sqlite.NewSolutionRepositoryTemp(getTestDbName(), getTestSolutionName())
	s.Require().Nil(err)
	rounds, err := sqlite.NewRoundRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	problems, err := sqlite.NewProblemRepositoryTemp(getTestDbName(), getTestProblemsName())
	s.Require().Nil(err)
	s.handler = handlers.StartReviewStage{
		Replier:             s.replier,
		Rounds:              &rounds,
		Solutions:           &solutions,
		Participants:        &participants,
		SolutionDistributor: &solutiondistributor.SolutionDistributor{},
		ReviewersCount:      2,
		Postman:             mocks.NewPostman(),
	}

	distributor := problemdistributor.NewSimpleDistributor(&problems, 3)
	_, err = mocks.GenReviewPendingRound(&rounds, &participants, &solutions, &problems,
		&distributor, 10, 3, []int{1, 2, 6})
	s.Require().Nil(err)
}

func (s *startReviewStageTs) TestMain() {
}

func TestStartReviewStage(t *testing.T) {
	suite.Run(t, &startReviewStageTs{})
}
