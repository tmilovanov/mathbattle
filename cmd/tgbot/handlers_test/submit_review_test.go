package handlerstest

import (
	"testing"

	"mathbattle/cmd/tgbot/handlers"
	mreplier "mathbattle/cmd/tgbot/replier"
	"mathbattle/mocks"
	mathbattle "mathbattle/models"
	"mathbattle/repository/sqlite"

	"github.com/stretchr/testify/suite"
)

type submitReviewTestSuite struct {
	suite.Suite

	handler   handlers.SubmitReview
	problems  mathbattle.ProblemRepository
	solutions mathbattle.SolutionRepository
	chatID    int64
}

func (s *submitReviewTestSuite) SetupTest() {
	var err error

	participants, err := sqlite.NewParticipantRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	rounds, err := sqlite.NewRoundRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	reviews, err := sqlite.NewReviewRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	problems, err := sqlite.NewProblemRepositoryTemp(getTestDbName(), getTestProblemsName())
	s.Require().Nil(err)
	s.problems = &problems
	solutions, err := sqlite.NewSolutionRepositoryTemp(getTestDbName(), getTestSolutionName())
	s.Require().Nil(err)
	s.solutions = &solutions

	s.handler = handlers.SubmitReview{
		Replier:      mreplier.RussianReplier{},
		Participants: &participants,
		Rounds:       &rounds,
		Reviews:      &reviews,
	}
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoRoundNoParticipant() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoRound() {
	_, err := s.handler.Participants.Store(mocks.GenParticipants(1, 11)[0])
	s.Require().Nil(err)
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoParticipant() {
	_, err := s.handler.Participants.Store(mocks.GenParticipants(1, 11)[0])
	s.Require().Nil(err)
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNotSolved() {
}

func (s *submitReviewTestSuite) TestGetCorrectKeyboard() {
}

func (s *submitReviewTestSuite) TestSetWrongSolutionNumber() {
}

func (s *submitReviewTestSuite) TestSendWrongFormatReview() {
}

func (s *submitReviewTestSuite) TestSendCorrectFormatReview() {
}

func TestSubmitReviewHandler(t *testing.T) {
	suite.Run(t, &submitReviewTestSuite{})
}
