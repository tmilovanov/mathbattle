package handlerstest

import (
	"fmt"
	"strconv"
	"testing"

	"mathbattle/cmd/tgbot/handlers"
	mreplier "mathbattle/cmd/tgbot/replier"
	"mathbattle/database/sqlite"
	mathbattle "mathbattle/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tb "gopkg.in/tucnak/telebot.v2"
)

type submitSolutionTestSuite struct {
	suite.Suite
	replier        mreplier.Replier
	handler        handlers.SubmitSolution
	chatID         int64
	req            *require.Assertions
	curRound       mathbattle.Round
	curParticipant mathbattle.Participant
}

func (s *submitSolutionTestSuite) SetupTest() {
	var err error

	s.req = require.New(s.T())
	s.replier = mreplier.RussianReplier{}

	participants, err := sqlite.NewSQLParticipantRepositoryTemp(getTestDbName())
	s.req.Nil(err)
	solutions, err := sqlite.NewSQLSolutionRepositoryTemp(getTestDbName(), getTestSolutionName())
	s.req.Nil(err)
	rounds, err := sqlite.NewSQLRoundRepositoryTemp(getTestDbName())
	s.req.Nil(err)
	problems, err := sqlite.NewSQLProblemRepositoryTemp(getTestDbName(), getTestProblemsName())
	s.req.Nil(err)
	s.handler = handlers.SubmitSolution{
		Replier:      s.replier,
		Participants: &participants,
		Rounds:       &rounds,
		Solutions:    &solutions,
	}

	// Store participant
	participant, err := participants.Store(mathbattle.Participant{
		TelegramID: strconv.FormatInt(s.chatID, 10),
		Grade:      11,
	})
	s.req.Nil(err)
	s.curParticipant = participant

	// Create problems
	problem1, err := problems.Store(mathbattle.Problem{
		MinGrade: 11,
		MaxGrade: 11,
		Content:  []byte("1234567890"),
	})
	s.req.Nil(err)

	problem2, err := problems.Store(mathbattle.Problem{
		MinGrade: 11,
		MaxGrade: 11,
		Content:  []byte("0987654321"),
	})
	s.req.Nil(err)

	// Create round
	round := mathbattle.NewRound()
	round.ProblemDistribution[participant.ID] = []string{problem1.ID, problem2.ID}
	round, err = rounds.Store(round)
	s.req.Nil(err)
	s.curRound = round
}

func (s *submitSolutionTestSuite) sendPhotos(photos []tb.Message) {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)

	testSequence := []reqRespSequence{
		{textReq(""), mathbattle.NewRespWithKeyboard(s.replier.SolutionExpectProblemNumber(), "1", "2"), 1},
		{textReq("1"), mathbattle.NewRespWithKeyboard(s.replier.SolutionExpectPart(), s.replier.SolutionFinishUploading()), 2},
	}

	for i, photo := range photos {
		testSequence = append(testSequence,
			reqRespSequence{photo, mathbattle.NewRespWithKeyboard(s.replier.SolutionPartUploaded(i+1), s.replier.SolutionFinishUploading()), 2})
	}

	testSequence = append(testSequence,
		reqRespSequence{textReq(s.replier.SolutionFinishUploading()), mathbattle.NewResp(s.replier.SolutionUploadSuccess(len(photos))), -1})

	sendReqExpectRespSequence(s.req, &s.handler, ctx, testSequence)
}

func (s *submitSolutionTestSuite) sendPhotosTestDatabase(images []mathbattle.Image) {
	photos := []tb.Message{}
	for i, part := range images {
		fakePhotoPath := fmt.Sprintf("fake/photo/p%d.%s", i, part.Extension)
		photos = append(photos, photoReq("", fakePhotoPath, part.Content))

	}

	s.sendPhotos(photos)

	problemID := s.curRound.ProblemDistribution[s.curParticipant.ID][0]
	solution, err := s.handler.Solutions.Find(s.curRound.ID, s.curParticipant.ID, problemID)
	s.req.Nil(err)
	s.req.Equal(solution.Parts, images)
}

func (s *submitSolutionTestSuite) TestUnregistered() {
}

func (s *submitSolutionTestSuite) TestRegisteredNoRound() {
}

func (s *submitSolutionTestSuite) TestSendWrongFormatSolution() {
}

func (s *submitSolutionTestSuite) TestSendNoSolution() {
	ctx := mathbattle.NewTelegramUserContext(s.chatID)
	sendReqExpectRespSequence(s.req, &s.handler, ctx, []reqRespSequence{
		{textReq(""), mathbattle.NewRespWithKeyboard(s.replier.SolutionExpectProblemNumber(), "1", "2"), 1},
		{textReq("1"), mathbattle.NewRespWithKeyboard(s.replier.SolutionExpectPart(), s.replier.SolutionFinishUploading()), 2},
		{textReq(s.replier.SolutionFinishUploading()), mathbattle.NewResp(s.replier.SolutionEmpty()), -1},
	})
}

func (s *submitSolutionTestSuite) TestSendOnePageSolution() {
	s.sendPhotosTestDatabase([]mathbattle.Image{
		{Extension: ".jpg", Content: []byte("1fakecontent")},
	})
}

func (s *submitSolutionTestSuite) TestSendMultiplePageSolution() {
	s.sendPhotosTestDatabase([]mathbattle.Image{
		{Extension: ".jpg", Content: []byte("1fakecontent")},
		{Extension: ".png", Content: []byte("2fakecontent")},
		{Extension: ".png", Content: []byte("3fakecontent")},
	})
}

func TestSubmitSolutionHandler(t *testing.T) {
	suite.Run(t, &submitSolutionTestSuite{})
}
