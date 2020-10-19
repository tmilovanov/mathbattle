package handlerstest

import (
	"strconv"
	"testing"

	"mathbattle/cmd/tgbot/handlers"
	"mathbattle/cmd/tgbot/replier"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"mathbattle/repository/sqlite"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type subscribeTestSuite struct {
	suite.Suite

	replier      mreplier.Replier
	participants sqlite.ParticipantRepository
	handler      handlers.Subscribe
	chatID       int64
	req          *require.Assertions
}

func (s *subscribeTestSuite) SetupTest() {
	s.req = require.New(s.T())
	s.replier = replier.RussianReplier{}
	participants, err := sqlite.NewParticipantRepositoryTemp(getTestDbName())
	s.participants = participants
	s.req.Nil(err)
	s.handler = handlers.Subscribe{
		Replier:      s.replier,
		Participants: &s.participants,
	}
}

func (s *subscribeTestSuite) TestCorrectSubscribe() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	testParticipant := mathbattle.Participant{
		TelegramID: strconv.FormatInt(s.chatID, 10),
		Name:       "JackDaniels",
		School:     "",
		Grade:      7,
	}

	sendTextExpectTextSequence(s.req, &s.handler, ctx, []reqRespTextSequence{
		{"", s.replier.RegisterNameExpect(), 1},
		{testParticipant.Name, s.replier.RegisterGradeExpect(), 2},
		{strconv.Itoa(testParticipant.Grade), s.replier.RegisterSuccess(), -1},
	})

	p, err := s.participants.GetByTelegramID(strconv.FormatInt(s.chatID, 10))
	testParticipant.ID = p.ID
	testParticipant.RegistrationTime = p.RegistrationTime
	s.req.Nil(err)
	s.req.Equal(p, testParticipant)
}

func (s *subscribeTestSuite) TestIncorrectName() {
	sendTextExpectTextSequence(s.req, &s.handler, mathbattle.NewTelegramUserContext(s.chatID), []reqRespTextSequence{
		{"", s.replier.RegisterNameExpect(), 1},
		{"123455~!!", s.replier.RegisterNameWrong(), 1},
		{"718317+-++", s.replier.RegisterNameWrong(), 1},
	})
}

func (s *subscribeTestSuite) TestIncorrectGrade() {
	sendTextExpectTextSequence(s.req, &s.handler, mathbattle.NewTelegramUserContext(s.chatID), []reqRespTextSequence{
		{"", s.replier.RegisterNameExpect(), 1},
		{"Jack", s.replier.RegisterGradeExpect(), 2},
		{"asdfsadf", s.replier.RegisterGradeWrong(), 2},
		{"-1", s.replier.RegisterGradeWrong(), 2},
		{"12", s.replier.RegisterGradeWrong(), 2},
	})
}

func (s *subscribeTestSuite) TestIncorrectThenCorrect() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	testParticipant := mathbattle.Participant{
		TelegramID: strconv.FormatInt(s.chatID, 10),
		Name:       "JackDaniels",
		School:     "",
		Grade:      7,
	}

	sendTextExpectTextSequence(s.req, &s.handler, ctx, []reqRespTextSequence{
		{"", s.replier.RegisterNameExpect(), 1},
		{"123455~!!", s.replier.RegisterNameWrong(), 1},
		{testParticipant.Name, s.replier.RegisterGradeExpect(), 2},
		{"12", s.replier.RegisterGradeWrong(), 2},
		{strconv.Itoa(testParticipant.Grade), s.replier.RegisterSuccess(), -1},
	})

	p, err := s.participants.GetByTelegramID(strconv.FormatInt(s.chatID, 10))
	testParticipant.ID = p.ID
	testParticipant.RegistrationTime = p.RegistrationTime
	s.req.Nil(err)
	s.req.Equal(p, testParticipant)
}

func (s *subscribeTestSuite) TestSubscirbeAlredyRegistered() {
	_, err := s.participants.Store(mathbattle.Participant{
		ID:         "",
		TelegramID: strconv.FormatInt(s.chatID, 10),
		Name:       "Jack",
		Grade:      7,
	})

	s.req.Nil(err)

	sendTextExpectTextSequence(s.req, &s.handler, mathbattle.NewTelegramUserContext(s.chatID), []reqRespTextSequence{
		{"", s.replier.AlreadyRegistered(), -1},
	})
}

func TestSubscribeHandler(t *testing.T) {
	suite.Run(t, &subscribeTestSuite{})
}
