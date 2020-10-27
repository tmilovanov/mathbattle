package handlerstest

import (
	"testing"

	"mathbattle/cmd/tgbot/handlers"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	"github.com/stretchr/testify/suite"
)

type submitReviewTestSuite struct {
	suite.Suite

	handler        handlers.SubmitReview
	replier        mreplier.Replier
	chatID         int64
	curRound       mathbattle.Round
	curParticipant mathbattle.Participant
}

func (s *submitReviewTestSuite) SetupTest() {
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoRoundNoParticipant() {
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoRound() {
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoParticipant() {
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
