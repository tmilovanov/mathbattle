package handlerstest

import (
	"strconv"
	"testing"

	"mathbattle/infrastructure"
	"mathbattle/interfaces/bot/handlers"
	"mathbattle/models/mathbattle"

	"github.com/stretchr/testify/suite"
)

type subscribeTs struct {
	suite.Suite

	handler handlers.Subscribe
	ctx     infrastructure.TelegramUserContext
}

func (s *subscribeTs) SetupTest() {
	container := infrastructure.NewTestContainer()

	s.ctx = infrastructure.NewTelegramUserContext(12345)
	s.handler = handlers.Subscribe{
		Replier:            container.Replier(),
		ParticipantService: container.ParticipantService(),
	}
}

func (s *subscribeTs) TestIncorrectName() {
	sendTextExpectTextSequence(s.Require(), &s.handler, s.ctx, []reqRespTextSequence{
		{"", s.handler.Replier.RegisterNameExpect(), 1},
		{"123455~!!", s.handler.Replier.RegisterNameWrong(), 1},
		{"718317+-++", s.handler.Replier.RegisterNameWrong(), 1},
	})
}

func (s *subscribeTs) TestIncorrectGrade() {
	sendTextExpectTextSequence(s.Require(), &s.handler, s.ctx, []reqRespTextSequence{
		{"", s.handler.Replier.RegisterNameExpect(), 1},
		{"Jack", s.handler.Replier.RegisterGradeExpect(), 2},
		{"asdfsadf", s.handler.Replier.RegisterGradeWrong(), 2},
		{"-1", s.handler.Replier.RegisterGradeWrong(), 2},
		{"12", s.handler.Replier.RegisterGradeWrong(), 2},
	})
}

func (s *subscribeTs) TestCorrectSubscribe() {
	testParticipant := mathbattle.Participant{
		TelegramID: s.ctx.User.ChatID,
		Name:       "JackDaniels",
		School:     "",
		Grade:      7,
	}

	sendTextExpectTextSequence(s.Require(), &s.handler, s.ctx, []reqRespTextSequence{
		{"", s.handler.Replier.RegisterNameExpect(), 1},
		{testParticipant.Name, s.handler.Replier.RegisterGradeExpect(), 2},
		{strconv.Itoa(testParticipant.Grade), s.handler.Replier.RegisterSuccess(), -1},
	})

	p, err := s.handler.ParticipantService.GetByTelegramID(s.ctx.User.ChatID)
	testParticipant.ID = p.ID
	testParticipant.RegistrationTime = p.RegistrationTime
	s.Require().Nil(err)
	s.Require().Equal(p, testParticipant)
}

func TestSubscribeHandler(t *testing.T) {
	suite.Run(t, &subscribeTs{})
}
