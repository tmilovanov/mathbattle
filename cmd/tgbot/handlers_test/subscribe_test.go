package handlers_test

import (
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

func tresp(input string) mathbattle.TelegramResponse {
	return mathbattle.TelegramResponse(input)
}

type MainTestSuite struct {
	suite.Suite
	replyer      mreplyer.Replier
	participants mock.MockParticipantsRepository
	handler      handlers.Subscribe
	chatID       int64
	req          *require.Assertions
}

type ReqRespTextSequence struct {
	request  string
	response string
	step     int
}

func (s *MainTestSuite) SetupTest() {
	s.replyer = replier.RussianReplyer{}
	s.participants = mock.NewMockParticipantsRepository()
	s.handler = handlers.Subscribe{
		Replier:      s.replyer,
		Participants: &s.participants,
	}
	s.req = require.New(s.T())
}

func (s *MainTestSuite) SendTextExpectText(ctx mathbattle.TelegramUserContext, msg string,
	expect string, expectedStep int) mathbattle.TelegramUserContext {

	result := ctx

	m := &tb.Message{Text: msg}
	step, resp, err := s.handler.Handle(ctx, m)
	s.req.Nil(err)
	s.req.Equal(resp, mathbattle.TelegramResponse(expect))
	s.req.Equal(step, expectedStep)
	result.CurrentStep = step

	return result
}

func (s *MainTestSuite) SendTextExpectTextSequence(ctx mathbattle.TelegramUserContext, seq []ReqRespTextSequence) {
	for _, elem := range seq {
		ctx = s.SendTextExpectText(ctx, elem.request, elem.response, elem.step)
	}
}

func (s *MainTestSuite) TestCorrectSubscribe() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	testParticipant := mathbattle.Participant{
		TelegramID: strconv.FormatInt(s.chatID, 10),
		Name:       "JackDaniels",
		School:     "",
		Grade:      7,
	}

	s.SendTextExpectTextSequence(ctx, []ReqRespTextSequence{
		{"", s.replyer.GetReply(mreplyer.ReplyRegisterNameExpect), 1},
		{testParticipant.Name, s.replyer.GetReply(mreplyer.ReplyRegisterGradeExpect), 2},
		{strconv.Itoa(testParticipant.Grade), s.replyer.GetReply(mreplyer.ReplyRegisterSuccess), -1},
	})

	p, err := s.participants.GetByTelegramID(strconv.FormatInt(s.chatID, 10))
	testParticipant.ID = p.ID
	testParticipant.RegistrationTime = p.RegistrationTime
	s.req.Nil(err)
	s.req.Equal(p, testParticipant)
}

func (s *MainTestSuite) TestIncorrectName() {
	s.SendTextExpectTextSequence(mathbattle.NewTelegramUserContext(s.chatID), []ReqRespTextSequence{
		{"", s.replyer.GetReply(mreplyer.ReplyRegisterNameExpect), 1},
		{"123455~!!", s.replyer.GetReply(mreplyer.ReplyRegisterNameWrong), 1},
		{"718317+-++", s.replyer.GetReply(mreplyer.ReplyRegisterNameWrong), 1},
	})
}

func (s *MainTestSuite) TestIncorrectGrade() {
	s.SendTextExpectTextSequence(mathbattle.NewTelegramUserContext(s.chatID), []ReqRespTextSequence{
		{"", s.replyer.GetReply(mreplyer.ReplyRegisterNameExpect), 1},
		{"Jack", s.replyer.GetReply(mreplyer.ReplyRegisterGradeExpect), 2},
		{"asdfsadf", s.replyer.GetReply(mreplyer.ReplyRegisterGradeWrong), 2},
		{"-1", s.replyer.GetReply(mreplyer.ReplyRegisterGradeWrong), 2},
		{"12", s.replyer.GetReply(mreplyer.ReplyRegisterGradeWrong), 2},
	})
}

func (s *MainTestSuite) TestIncorrectThenCorrect() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	testParticipant := mathbattle.Participant{
		TelegramID: strconv.FormatInt(s.chatID, 10),
		Name:       "JackDaniels",
		School:     "",
		Grade:      7,
	}

	s.SendTextExpectTextSequence(ctx, []ReqRespTextSequence{
		{"", s.replyer.GetReply(mreplyer.ReplyRegisterNameExpect), 1},
		{"123455~!!", s.replyer.GetReply(mreplyer.ReplyRegisterNameWrong), 1},
		{testParticipant.Name, s.replyer.GetReply(mreplyer.ReplyRegisterGradeExpect), 2},
		{"12", s.replyer.GetReply(mreplyer.ReplyRegisterGradeWrong), 2},
		{strconv.Itoa(testParticipant.Grade), s.replyer.GetReply(mreplyer.ReplyRegisterSuccess), -1},
	})

	p, err := s.participants.GetByTelegramID(strconv.FormatInt(s.chatID, 10))
	testParticipant.ID = p.ID
	testParticipant.RegistrationTime = p.RegistrationTime
	s.req.Nil(err)
	s.req.Equal(p, testParticipant)
}

func (s *MainTestSuite) TestSubscirbeAlredyRegistered() {
	s.participants.Store(mathbattle.Participant{
		ID:         "",
		TelegramID: strconv.FormatInt(s.chatID, 10),
		Name:       "Jack",
		Grade:      7,
	})

	s.SendTextExpectTextSequence(mathbattle.NewTelegramUserContext(s.chatID), []ReqRespTextSequence{
		{"", s.replyer.GetReply(mreplyer.ReplyAlreadyRegistered), -1},
	})
}

func TestSubscribeHandler(t *testing.T) {
	suite.Run(t, &MainTestSuite{})
}
