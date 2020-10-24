package handlerstest

import (
	"strconv"
	"testing"

	"mathbattle/cmd/tgbot/handlers"
	"mathbattle/cmd/tgbot/replier"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"mathbattle/repository/sqlite"

	"github.com/stretchr/testify/suite"
)

type subscribeTs struct {
	suite.Suite

	handler      handlers.Subscribe
	replier      mreplier.Replier
	participants sqlite.ParticipantRepository
	chatID       int64
}

func (s *subscribeTs) SetupTest() {
	s.replier = replier.RussianReplier{}
	participants, err := sqlite.NewParticipantRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	s.participants = participants
	s.handler = handlers.Subscribe{
		Replier:      s.replier,
		Participants: &s.participants,
	}
}

func (s *subscribeTs) TestCorrectSubscribe() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	testParticipant := mathbattle.Participant{
		TelegramID: s.chatID,
		Name:       "JackDaniels",
		School:     "",
		Grade:      7,
	}

	sendTextExpectTextSequence(s.Require(), &s.handler, ctx, []reqRespTextSequence{
		{"", s.replier.RegisterNameExpect(), 1},
		{testParticipant.Name, s.replier.RegisterGradeExpect(), 2},
		{strconv.Itoa(testParticipant.Grade), s.replier.RegisterSuccess(), -1},
	})

	p, err := s.participants.GetByTelegramID(s.chatID)
	testParticipant.ID = p.ID
	testParticipant.RegistrationTime = p.RegistrationTime
	s.Require().Nil(err)
	s.Require().Equal(p, testParticipant)
}

func (s *subscribeTs) TestIncorrectName() {
	sendTextExpectTextSequence(s.Require(), &s.handler, mathbattle.NewTelegramUserContext(s.chatID), []reqRespTextSequence{
		{"", s.replier.RegisterNameExpect(), 1},
		{"123455~!!", s.replier.RegisterNameWrong(), 1},
		{"718317+-++", s.replier.RegisterNameWrong(), 1},
	})
}

func (s *subscribeTs) TestIncorrectGrade() {
	sendTextExpectTextSequence(s.Require(), &s.handler, mathbattle.NewTelegramUserContext(s.chatID), []reqRespTextSequence{
		{"", s.replier.RegisterNameExpect(), 1},
		{"Jack", s.replier.RegisterGradeExpect(), 2},
		{"asdfsadf", s.replier.RegisterGradeWrong(), 2},
		{"-1", s.replier.RegisterGradeWrong(), 2},
		{"12", s.replier.RegisterGradeWrong(), 2},
	})
}

func (s *subscribeTs) TestIncorrectThenCorrect() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	testParticipant := mathbattle.Participant{
		TelegramID: s.chatID,
		Name:       "JackDaniels",
		School:     "",
		Grade:      7,
	}

	sendTextExpectTextSequence(s.Require(), &s.handler, ctx, []reqRespTextSequence{
		{"", s.replier.RegisterNameExpect(), 1},
		{"123455~!!", s.replier.RegisterNameWrong(), 1},
		{testParticipant.Name, s.replier.RegisterGradeExpect(), 2},
		{"12", s.replier.RegisterGradeWrong(), 2},
		{strconv.Itoa(testParticipant.Grade), s.replier.RegisterSuccess(), -1},
	})

	p, err := s.participants.GetByTelegramID(s.chatID)
	testParticipant.ID = p.ID
	testParticipant.RegistrationTime = p.RegistrationTime
	s.Require().Nil(err)
	s.Require().Equal(p, testParticipant)
}

func TestSubscribeHandler(t *testing.T) {
	suite.Run(t, &subscribeTs{})
}
