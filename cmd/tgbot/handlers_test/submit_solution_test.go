package handlers_test

import (
	"bytes"
	"strconv"
	"testing"

	"mathbattle/cmd/tgbot/handlers"
	"mathbattle/cmd/tgbot/replier"
	mreplyer "mathbattle/cmd/tgbot/replier"
	"mathbattle/database/mock"
	mathbattle "mathbattle/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tb "gopkg.in/tucnak/telebot.v2"
)

type submitSolutionTestSuite struct {
	suite.Suite
	replyer      mreplyer.Replier
	participants mock.MockParticipantsRepository
	solutions    mock.SolutionRepository
	rounds       mock.RoundRepository
	handler      handlers.SubmitSolution
	chatID       int64
	req          *require.Assertions
}

func (s *submitSolutionTestSuite) GetFakeFile(fakePath string, fakeContent []byte) tb.File {
	return tb.File{
		FilePath:   fakePath,
		FileReader: bytes.NewReader(fakeContent),
	}
}

func (s *submitSolutionTestSuite) SetupTest() {
	s.replyer = replier.RussianReplyer{}
	s.participants = mock.NewMockParticipantsRepository()
	s.solutions = mock.NewSolutionRepository()
	s.rounds = mock.NewRoundRepository()
	s.handler = handlers.SubmitSolution{
		Replier:      s.replyer,
		Participants: &s.participants,
		Rounds:       &s.rounds,
		Solutions:    &s.solutions,
	}
	s.req = require.New(s.T())
}

func (s *submitSolutionTestSuite) TestUnregistered() {
}

func (s *submitSolutionTestSuite) TestRegisteredNoRound() {
}

func (s *submitSolutionTestSuite) TestSendWrongFormatSolution() {
}

func (s *submitSolutionTestSuite) TestSendOnePageSolution() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)

	s.participants.Store(mathbattle.Participant{TelegramID: strconv.FormatInt(ctx.ChatID, 10)})
	mathbattle.NewRound()
	s.rounds.Store()

	msg := tb.Message{Text: ""}
	step, resp, err := s.handler.Handle(ctx, &msg)
	s.req.Nil(err)
	s.req.Equal(resp, "")
	s.req.Equal(step, 1)

	//fakePath := "page1.jpg"
	//fakeContent := []byte("1234567890")
	//msg = tb.Message{
	//Photo: &tb.Photo{
	//File: tb.File{
	//FilePath:   fakePath,
	//FileReader: bytes.NewReader(fakeContent),
	//},
	//},
	//}
	//step, resp, err := s.handler.Handle(ctx, &msg)
	//s.req.Nil(err)
}

func (s *submitSolutionTestSuite) TestSendMultiplePageSolution() {
}

func TestSubmitSolutionHandler(t *testing.T) {
	suite.Run(t, &submitSolutionTestSuite{})
}
