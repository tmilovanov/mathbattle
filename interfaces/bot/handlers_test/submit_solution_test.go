package handlerstest

import (
	"fmt"
	"testing"

	"mathbattle/infrastructure"
	"mathbattle/interfaces/bot/handlers"
	"mathbattle/models/mathbattle"

	"github.com/stretchr/testify/suite"
	tb "gopkg.in/tucnak/telebot.v2"
)

type submitSolutionTestSuite struct {
	suite.Suite

	handler        handlers.SubmitSolution
	chatID         int64
	curRound       mathbattle.Round
	curParticipant mathbattle.Participant
}

func (s *submitSolutionTestSuite) SetupTest() {
	container := infrastructure.NewTestContainer()

	s.handler = handlers.SubmitSolution{
		Replier:            container.Replier(),
		ParticipantService: container.ParticipantService(),
		RoundService:       container.RoundService(),
		SolutionService:    container.SolutionService(),
	}

	s.curRound = container.CreateSolveStageRound(infrastructure.TestRoundDescription{
		ParticipantsCount: 1,
		ProblemsOnEach:    2,
	})
}

func (s *submitSolutionTestSuite) sendSolutionParts(ctx infrastructure.TelegramUserContext,
	images []mathbattle.Image) infrastructure.TelegramUserContext {

	testSequence := []reqRespSequence{}
	for i, part := range images {
		fakePhotoPath := fmt.Sprintf("fake/photo/p%d.%s", i, part.Extension)
		testSequence = append(testSequence, reqRespSequence{
			request:  photo("", fakePhotoPath, part.Content),
			response: handlers.NewRespWithKeyboard(s.handler.Replier.SolutionPartUploaded(i+1), s.handler.Replier.SolutionFinishUploading()),
			step:     3,
		})
	}

	return sendReqExpectRespSequence(s.Require(), &s.handler, ctx, testSequence)
}

func (s *submitSolutionTestSuite) sendPhotos(photos []tb.Message) infrastructure.TelegramUserContext {
	ctx := infrastructure.NewTelegramUserContext(s.chatID)

	testSequence := []reqRespSequence{
		{text(""), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectProblemCaption(), "A", "B"), 1},
		{text("A"), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectPart(), s.handler.Replier.SolutionFinishUploading()), 3},
	}

	for i, photo := range photos {
		testSequence = append(testSequence,
			reqRespSequence{photo, handlers.NewRespWithKeyboard(s.handler.Replier.SolutionPartUploaded(i+1), s.handler.Replier.SolutionFinishUploading()), 3})
	}

	testSequence = append(testSequence,
		reqRespSequence{text(s.handler.Replier.SolutionFinishUploading()), handlers.NewResp(s.handler.Replier.SolutionUploadSuccess(len(photos))), -1})

	return sendReqExpectRespSequence(s.Require(), &s.handler, ctx, testSequence)
}

func (s *submitSolutionTestSuite) sendPhotosTestDatabase(images []mathbattle.Image) infrastructure.TelegramUserContext {
	photos := []tb.Message{}
	for i, part := range images {
		fakePhotoPath := fmt.Sprintf("fake/photo/p%d.%s", i, part.Extension)
		photos = append(photos, photo("", fakePhotoPath, part.Content))

	}

	ctx := s.sendPhotos(photos)

	problemID := s.curRound.ProblemDistribution[s.curParticipant.ID][0].ProblemID
	solutions, err := s.handler.SolutionService.Find(mathbattle.FindDescriptor{
		RoundID:       s.curRound.ID,
		ParticipantID: s.curParticipant.ID,
		ProblemID:     problemID,
	})
	s.Require().Nil(err)
	solution := solutions[0]

	s.Require().Equal(len(images), len(solution.Parts)) // Optional, but easy to see the difference in console
	s.Require().Equal(images, solution.Parts)

	return ctx
}

func (s *submitSolutionTestSuite) TestSendNoSolution() {
	ctx := infrastructure.NewTelegramUserContext(s.chatID)
	sendReqExpectRespSequence(s.Require(), &s.handler, ctx, []reqRespSequence{
		{text(""), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectProblemCaption(), "A", "B"), 1},
		{text("A"), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectPart(), s.handler.Replier.SolutionFinishUploading()), 3},
		{text(s.handler.Replier.SolutionFinishUploading()), handlers.NewResp(s.handler.Replier.SolutionEmpty()), -1},
	})
}

func (s *submitSolutionTestSuite) TestSendSolutionFirstTime() {
	ctx := infrastructure.NewTelegramUserContext(s.chatID)
	sendReqExpectRespSequence(s.Require(), &s.handler, ctx, []reqRespSequence{
		{text(""), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectProblemCaption(), "A", "B"), 1},
		{text("A"), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectPart(), s.handler.Replier.SolutionFinishUploading()), 3},
		{text("BlahBlah"), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionWrongFormat(), s.handler.Replier.SolutionFinishUploading()), 3},
		{photo("", "fake/path/p1.jpg", []byte("12345")),
			handlers.NewRespWithKeyboard(s.handler.Replier.SolutionPartUploaded(1), s.handler.Replier.SolutionFinishUploading()), 3},
		{text("AsdfAsdf"), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionWrongFormat(), s.handler.Replier.SolutionFinishUploading()), 3},
		{photo("", "fake/path/p2.jpg", []byte("54321")),
			handlers.NewRespWithKeyboard(s.handler.Replier.SolutionPartUploaded(2), s.handler.Replier.SolutionFinishUploading()), 3},
		{text(s.handler.Replier.SolutionFinishUploading()), handlers.NewResp(s.handler.Replier.SolutionUploadSuccess(2)), -1},
	})
}

func (s *submitSolutionTestSuite) TestSendSolutionSecondTime() {
	ctx := infrastructure.NewTelegramUserContext(s.chatID)
	sendReqExpectRespSequence(s.Require(), &s.handler, ctx, []reqRespSequence{
		{text(""), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectProblemCaption(), "A", "B"), 1},
		{text("A"), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectPart(), s.handler.Replier.SolutionFinishUploading()), 3},
		{photo("", "fake/path/p1.jpg", []byte("12345")),
			handlers.NewRespWithKeyboard(s.handler.Replier.SolutionPartUploaded(1), s.handler.Replier.SolutionFinishUploading()), 3},
		{photo("", "fake/path/p2.jpg", []byte("54321")),
			handlers.NewRespWithKeyboard(s.handler.Replier.SolutionPartUploaded(2), s.handler.Replier.SolutionFinishUploading()), 3},
		{text(s.handler.Replier.SolutionFinishUploading()), handlers.NewResp(s.handler.Replier.SolutionUploadSuccess(2)), -1},
		// Send again, but refuse to rewrite
		{text(""), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectProblemCaption(), "A", "B"), 1},
		{text("A"), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionIsRewriteOld(), s.handler.Replier.Yes(), s.handler.Replier.No()), 2},
		{text(s.handler.Replier.No()), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionDeclineRewriteOld()), -1},
		// Send again, but agree to rewrite
		{text(""), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectProblemCaption(), "A", "B"), 1},
		{text("A"), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionIsRewriteOld(), s.handler.Replier.Yes(), s.handler.Replier.No()), 2},
		{text(s.handler.Replier.Yes()), handlers.NewRespWithKeyboard(s.handler.Replier.SolutionExpectPart(), s.handler.Replier.SolutionFinishUploading()), 3},
		{photo("", "fake/path/p1.jpg", []byte("12345")),
			handlers.NewRespWithKeyboard(s.handler.Replier.SolutionPartUploaded(1), s.handler.Replier.SolutionFinishUploading()), 3},
		{photo("", "fake/path/p2.jpg", []byte("54321")),
			handlers.NewRespWithKeyboard(s.handler.Replier.SolutionPartUploaded(2), s.handler.Replier.SolutionFinishUploading()), 3},
		{text(s.handler.Replier.SolutionFinishUploading()), handlers.NewResp(s.handler.Replier.SolutionUploadSuccess(2)), -1},
	})
}

func TestSubmitSolutionHandler(t *testing.T) {
	suite.Run(t, &submitSolutionTestSuite{})
}
