package handlerstest

import (
	"fmt"
	"strconv"
	"testing"

	"mathbattle/cmd/tgbot/handlers"
	mreplier "mathbattle/cmd/tgbot/replier"
	"mathbattle/mocks"
	mathbattle "mathbattle/models"
	problemdistributor "mathbattle/problem_distributor"
	"mathbattle/repository/sqlite"
	solutiondistributor "mathbattle/solution_distributor"

	"github.com/stretchr/testify/suite"
)

func findParticipantWithoutReview(round mathbattle.Round,
	participants mathbattle.ParticipantRepository) (mathbattle.Participant, error) {

	allParticipants, err := participants.GetAll()
	if err != nil {
		return mathbattle.Participant{}, err
	}

	for _, participant := range allParticipants {
		solutionsIDs := round.ReviewDistribution.BetweenParticipants[participant.ID]
		if len(solutionsIDs) == 0 {
			return participant, nil
		}
	}

	return mathbattle.Participant{}, mathbattle.ErrNotFound
}

func findParticipantWithReview(round mathbattle.Round,
	participants mathbattle.ParticipantRepository) (mathbattle.Participant, error) {

	allParticipants, err := participants.GetAll()
	if err != nil {
		return mathbattle.Participant{}, err
	}

	for _, participant := range allParticipants {
		solutionsIDs := round.ReviewDistribution.BetweenParticipants[participant.ID]
		if len(solutionsIDs) != 0 {
			return participant, nil
		}
	}

	return mathbattle.Participant{}, mathbattle.ErrNotFound
}

type submitReviewTestSuite struct {
	suite.Suite

	handler   handlers.SubmitReview
	problems  mathbattle.ProblemRepository
	solutions mathbattle.SolutionRepository
}

func (s *submitReviewTestSuite) SetupTest() {
	var err error

	participants, err := sqlite.NewParticipantRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	rounds, err := sqlite.NewRoundRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	reviews, err := sqlite.NewReviewRepositoryTemp(getTestDbName())
	s.Require().Nil(err)
	problems, err := sqlite.NewProblemRepositoryTemp(getTestDbName(), getTestProblemsName())
	s.Require().Nil(err)
	s.problems = &problems
	solutions, err := sqlite.NewSolutionRepositoryTemp(getTestDbName(), getTestSolutionName())
	s.Require().Nil(err)
	s.solutions = &solutions

	s.handler = handlers.SubmitReview{
		Replier:      mreplier.RussianReplier{},
		Participants: &participants,
		Rounds:       &rounds,
		Reviews:      &reviews,
		Solutions:    s.solutions,
		Postman:      mocks.NewPostman(),
	}
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoRoundNoParticipant() {
	ctx := mathbattle.NewTelegramUserContext(12345)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoRound() {
	_, err := s.handler.Participants.Store(mocks.GenParticipants(1, 11)[0])
	s.Require().Nil(err)
	ctx := mathbattle.NewTelegramUserContext(12345)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoParticipant() {
	// Generate round in review running stage where current telegram user is not a participant
	problemDistributor := problemdistributor.NewSimpleDistributor(s.problems, 3)
	_, err := mocks.GenReviewStageRound(s.handler.Rounds, s.handler.Participants, s.solutions,
		s.problems, &problemDistributor, &solutiondistributor.SolutionDistributor{},
		10, 3, []int{2, 2, 2}, 2)
	s.Require().Nil(err)
	ctx := mathbattle.NewTelegramUserContext(12345)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *submitReviewTestSuite) TestNotSuitableWhenNoSolutionWasSent() {
	// Generate round in review running stage where current telegram user is participant,
	// but didn't get any solution on review
	problemDistributor := problemdistributor.NewSimpleDistributor(s.problems, 3)
	round, err := mocks.GenReviewStageRound(s.handler.Rounds, s.handler.Participants, s.solutions,
		s.problems, &problemDistributor, &solutiondistributor.SolutionDistributor{},
		10, 3, []int{2, 2, 2}, 2)
	s.Require().Nil(err)

	p, err := findParticipantWithoutReview(round, s.handler.Participants)
	s.Require().Nil(err)

	ctx := mathbattle.NewTelegramUserContext(p.TelegramID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().False(isSuitable)
}

func (s *submitReviewTestSuite) TestSuitable() {
	// Generate round in review running stage where current telegram user is participant,
	// and got some solutions on review
	problemDistributor := problemdistributor.NewSimpleDistributor(s.problems, 3)
	round, err := mocks.GenReviewStageRound(s.handler.Rounds, s.handler.Participants, s.solutions,
		s.problems, &problemDistributor, &solutiondistributor.SolutionDistributor{},
		10, 3, []int{2, 2, 2}, 2)
	s.Require().Nil(err)

	p, err := findParticipantWithReview(round, s.handler.Participants)
	s.Require().Nil(err)

	ctx := mathbattle.NewTelegramUserContext(p.TelegramID)
	isSuitable, err := s.handler.IsCommandSuitable(ctx)
	s.Require().Nil(err)
	s.Require().True(isSuitable)
}

func (s *submitReviewTestSuite) TestSetOverwriteReview() {
	// Set review
	// Decline overwrite
	// Overwrite
	problemDistributor := problemdistributor.NewSimpleDistributor(s.problems, 3)
	round, err := mocks.GenReviewStageRound(s.handler.Rounds, s.handler.Participants, s.solutions,
		s.problems, &problemDistributor, &solutiondistributor.SolutionDistributor{},
		10, 3, []int{10, 10, 10}, 2)
	s.Require().Nil(err)

	p, err := findParticipantWithReview(round, s.handler.Participants)
	s.Require().Nil(err)

	ctx := mathbattle.NewTelegramUserContext(p.TelegramID)
	solutionsNumbers := mathbattle.SolutionNumbers(round, p)
	for _, sn := range solutionsNumbers {
		solutionI, _ := strconv.Atoi(sn)
		solutionI--
		solutionID := round.ReviewDistribution.BetweenParticipants[p.ID][solutionI]

		reviewContent := fmt.Sprintf("example_review_%s", sn)
		sendReqExpectRespSequence(s.Require(), &s.handler, ctx, []reqRespSequence{
			{textReq(""), mathbattle.NewRespWithKeyboard(s.handler.Replier.ReviewExpectSolutionNumber(), solutionsNumbers...), 1},
			{textReq(sn), mathbattle.NewResp(s.handler.Replier.ReviewExpectContent()), 3},
			{textReq(reviewContent), mathbattle.NewResp(s.handler.Replier.ReviewUploadSuccess()), -1},
		})

		reviews, err := s.handler.Reviews.FindMany(p.ID, solutionID)
		s.Require().Nil(err)
		s.Require().Equal(1, len(reviews))
		s.Require().Equal(reviewContent, reviews[0].Content)

		sendReqExpectRespSequence(s.Require(), &s.handler, ctx, []reqRespSequence{
			{textReq(""), mathbattle.NewRespWithKeyboard(s.handler.Replier.ReviewExpectSolutionNumber(), solutionsNumbers...), 1},
			{textReq(sn), mathbattle.NewRespWithKeyboard(s.handler.Replier.ReviewIsRewriteOld(),
				s.handler.Replier.Yes(), s.handler.Replier.No()), 2},
			{textReq(s.handler.Replier.No()), mathbattle.NewResp(s.handler.Replier.Cancel()), -1},
		})

		reviewContent = fmt.Sprintf("example_review_overwrite_%s", sn)
		sendReqExpectRespSequence(s.Require(), &s.handler, ctx, []reqRespSequence{
			{textReq(""), mathbattle.NewRespWithKeyboard(s.handler.Replier.ReviewExpectSolutionNumber(), solutionsNumbers...), 1},
			{textReq(sn), mathbattle.NewRespWithKeyboard(s.handler.Replier.ReviewIsRewriteOld(),
				s.handler.Replier.Yes(), s.handler.Replier.No()), 2},
			{textReq(s.handler.Replier.Yes()), mathbattle.NewResp(s.handler.Replier.ReviewExpectContent()), 3},
			{textReq(reviewContent), mathbattle.NewResp(s.handler.Replier.ReviewUploadSuccess()), -1},
		})

		reviews, err = s.handler.Reviews.FindMany(p.ID, solutionID)
		s.Require().Nil(err)
		s.Require().Equal(1, len(reviews))
		s.Require().Equal(reviewContent, reviews[0].Content)
	}

}

func (s *submitReviewTestSuite) TestSetWrongSolutionNumber() {
	problemDistributor := problemdistributor.NewSimpleDistributor(s.problems, 3)
	round, err := mocks.GenReviewStageRound(s.handler.Rounds, s.handler.Participants, s.solutions,
		s.problems, &problemDistributor, &solutiondistributor.SolutionDistributor{},
		10, 3, []int{10, 10, 10}, 2)
	s.Require().Nil(err)

	p, err := findParticipantWithReview(round, s.handler.Participants)
	s.Require().Nil(err)

	ctx := mathbattle.NewTelegramUserContext(p.TelegramID)
	solutionsNumbers := mathbattle.SolutionNumbers(round, p)
	sendReqExpectRespSequence(s.Require(), &s.handler, ctx, []reqRespSequence{
		{textReq(""), mathbattle.NewRespWithKeyboard(s.handler.Replier.ReviewExpectSolutionNumber(), solutionsNumbers...), 1},
		{textReq("100"), mathbattle.NewRespWithKeyboard(s.handler.Replier.ReviewWrongSolutionNumber(), solutionsNumbers...), 1},
	})
}

func TestSubmitReviewHandler(t *testing.T) {
	suite.Run(t, &submitReviewTestSuite{})
}
