package handlerstest

import (
	"testing"
	"time"

	"mathbattle/cmd/tgbot/handlers"
	"mathbattle/cmd/tgbot/replier"
	mreplier "mathbattle/cmd/tgbot/replier"
	"mathbattle/mocks"
	mathbattle "mathbattle/models"
	"mathbattle/mstd"
	problemdistributor "mathbattle/problem_distributor"
	"mathbattle/repository/sqlite"

	"github.com/stretchr/testify/suite"
)

type getProblemsTs struct {
	suite.Suite

	handler      handlers.GetProblems
	chatID       int64
	replier      mreplier.Replier
	participants sqlite.ParticipantRepository
	rounds       sqlite.RoundRepository
	problems     sqlite.ProblemRepository
}

func (s *getProblemsTs) SetupTest() {
	s.chatID = 123456
	s.replier = replier.RussianReplier{}

	var err error
	s.participants, err = sqlite.NewParticipantRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	s.rounds, err = sqlite.NewRoundRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	s.problems, err = sqlite.NewProblemRepositoryTemp(getTestDbName(), getTestProblemsName())
	s.Require().Nil(err)

	s.handler = handlers.GetProblems{
		Participants: &s.participants,
		Rounds:       &s.rounds,
		Problems:     &s.problems,
	}
}

func (s *getProblemsTs) TestNotSuitableWhenNoRoundNoParticipant() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *getProblemsTs) TestNotSuitableWhenNoRound() {
	_, err := s.participants.Store(mathbattle.Participant{
		TelegramID: s.chatID,
		Name:       "Test Name",
	})
	s.Require().Nil(err)

	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *getProblemsTs) TestNotSuitableWhenNoParticipant() {
	duration, _ := time.ParseDuration("48h")
	_, err := s.rounds.Store(mathbattle.NewRound(duration))
	s.Require().Nil(err)

	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *getProblemsTs) TestSuitableWhenRoundAndParticipant() {
	_, err := s.participants.Store(mathbattle.Participant{
		TelegramID: s.chatID,
		Name:       "Test Name",
	})
	s.Require().Nil(err)

	duration, _ := time.ParseDuration("48h")
	_, err = s.rounds.Store(mathbattle.NewRound(duration))
	s.Require().Nil(err)

	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().True(isSuitable)
}

func (s *getProblemsTs) TestGetProblems() {
	distributor := problemdistributor.NewSimpleDistributor(&s.problems, 3)
	_, err := mocks.GenSolutionStageRound(&s.rounds, &s.participants, &s.problems, &distributor, 1, 3)
	s.Require().Nil(err)

	participants, err := s.participants.GetAll()
	s.Require().Nil(err)
	s.Require().Equal(1, len(participants))

	participants[0].TelegramID = s.chatID
	err = s.participants.Update(participants[0])
	s.Require().Nil(err)

	problems, err := s.problems.GetAll()
	s.Require().Nil(err)
	s.Require().Equal(3, len(problems))

	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	msg := textReq("")
	step, resp, err := s.handler.Handle(ctx, &msg)
	s.Require().Nil(err)
	s.Require().Equal(step, -1)
	s.Require().Equal(len(problems), len(resp))
	for i := 0; i < len(problems); i++ {
		problemImg := mathbattle.Image{
			Extension: problems[i].Extension,
			Content:   problems[i].Content,
		}
		s.Require().Equal(mstd.IndexToLetter(i), resp[i].Text)
		s.Require().Equal(problemImg, resp[i].Img)
	}
}

func TestGetProblemsHandler(t *testing.T) {
	suite.Run(t, &getProblemsTs{})
}
