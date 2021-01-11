package handlerstest

import (
	"strconv"
	"testing"

	"mathbattle/infrastructure"
	"mathbattle/infrastructure/repository/memory"
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

	ctxRepository, err := memory.NewTelegramContextRepository(container.UserRepository())
	s.Require().Nil(err)

	s.ctx, err = ctxRepository.GetByUserData(infrastructure.TelegramUserData{
		ChatID:   12345,
		Username: "FakeUser",
	})
	s.Require().Nil(err)

	s.handler = handlers.Subscribe{
		Replier:            container.Replier(),
		ParticipantService: container.ParticipantService(),
	}
}

func (s *subscribeTs) TestSubscribeNew() {
	testParticipant := mathbattle.Participant{
		Name:     "JackDaniels",
		School:   "",
		Grade:    7,
		IsActive: true,
	}

	sendTextExpectTextSequence(s.Require(), &s.handler, s.ctx, []reqRespTextSequence{
		{"", s.handler.Replier.RegisterNameExpect(), 1},
		// Try incorrect name
		{"123455~!!", s.handler.Replier.RegisterNameWrong(), 1},
		{"718317+-++", s.handler.Replier.RegisterNameWrong(), 1},
		// Correct name
		{testParticipant.Name, s.handler.Replier.RegisterGradeExpect(), 2},
		// Try incorrect grade
		{"Jack", s.handler.Replier.RegisterGradeWrong(), 2},
		{"asdfsadf", s.handler.Replier.RegisterGradeWrong(), 2},
		{"-1", s.handler.Replier.RegisterGradeWrong(), 2},
		{"12", s.handler.Replier.RegisterGradeWrong(), 2},
		// Correct grade
		{strconv.Itoa(testParticipant.Grade), s.handler.Replier.RegisterSuccess(), -1},
	})

	p, err := s.handler.ParticipantService.GetByTelegramID(s.ctx.User.TelegramID)
	testParticipant.User = s.ctx.User
	testParticipant.ID = p.ID
	s.Require().Nil(err)
	s.Require().Equal(p, testParticipant)
}

func (s *subscribeTs) TestSubscribeInactive() {
}

func TestSubscribeHandler(t *testing.T) {
	suite.Run(t, &subscribeTs{})
}
