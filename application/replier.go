package application

import (
	"time"

	"mathbattle/models/mathbattle"
)

type TelegramCommandHelp struct {
	Name string
	Desc string
}

type Replier interface {
	GetStartMessage() string
	GetAvailableCommands(availableCommands []TelegramCommandHelp) string
	GetHelpMessages() []string

	Yes() string
	No() string
	Cancel() string

	CmdHelpName() string
	CmdHelpDesc() string
	CmdSubscribeName() string
	CmdSubscribeDesc() string
	CmdUnsubscribeName() string
	CmdUnsubscribeDesc() string
	CmdGetProblemsName() string
	CmdGetProblemsDesc() string
	CmdSubmitSolutionName() string
	CmdSubmitSolutionDesc() string
	CmdStartReviewStageName() string
	CmdStartReviewStageDesc() string
	CmdSubmitReviewName() string
	CmdSubmitReviewDesc() string
	CmdStatName() string
	CmdStatDesc() string
	CmdGetReviewsName() string
	CmdGetReviewsDesc() string
	CmdStartRoundName() string
	CmdStartRoundDesc() string
	CmdServiceMsgName() string
	CmdServiceMsgDesc() string
	CmdGetMyResultsName() string
	CmdGetMyResultsDesc() string

	InternalError() string
	NotParticipant() string
	NoRoundRunning() string

	SolveStageEnd() string
	SolveStageEndNoSolutions() string
	ReviewStageEnd() string

	// Replies used in CmdSubscribe
	AlreadyRegistered() string
	RegisterNameExpect() string
	RegisterNameWrong() string
	RegisterGradeExpect() string
	RegisterGradeWrong() string
	RegisterSuccess() string
	RegisterSuccessRoundRunning(solveStageDuration time.Duration, solveStageEndMsk time.Time) string

	// Replies used in CmdUnsubscribe
	NotSubscribed() string
	UnsubscribeSuccess() string

	// Replies used to post problems during start of round
	ProblemsPostBefore(stageDuration time.Duration, stageEnd time.Time) string
	ProblemsPostAfter() string

	// Replies used in CmdSubmitSolution
	SolutionUploadSuccess(totalUpload int) string
	SolutionPartUploaded(partNumber int) string
	SolutionExpectProblemCaption() string
	SolutionWrongProblemCaptionFormat() string
	SolutionWrongProblemCaption() string
	SolutionExpectPart() string
	SolutionIsRewriteOld() string
	SolutionDeclineRewriteOld() string
	SolutionFinishUploading() string
	SolutionWrongFormat() string
	SolutionEmpty() string

	// Replies used in CmdStartReviewStage
	StartReviewGetDuration() string
	StartReviewWrongDuration() string
	StartReviewConfirmDuration(untilDate time.Time) string
	StartReviewSuccess() string

	// Replies used in CmdStartRound
	StartRoundGetDuration() string
	StartRoundWrongDuration() string
	StartRoundConfirmDuration(untillDate time.Time) string
	StartRoundSuccess(startResult mathbattle.StartResult) string
	StartRoundAskProblemsIDs() string

	// Replies used to post solutions to other participants to review
	ReviewPostBefore(stageDuration time.Duration, stageEnd time.Time) string
	ReviewPostCaption(problemCaption string, solutionNumber int) string
	ReviewPostAfter() string

	// Replies used in CmdSubmitReview
	ReviewGetSolutionCaptions(descriptors []mathbattle.SolutionDescriptor) []string
	ReviewGetDescriptor(userInput string) (mathbattle.SolutionDescriptor, bool)
	ReviewExpectSolutionCaption() string
	ReviewWrongSolutionCaption() string
	ReviewIsRewriteOld() string
	ReviewExpectContent() string
	ReviewUploadSuccess() string
	ReviewMsgForReviewee(review mathbattle.Review) string

	// Replies used in CmdStat
	FormatStat(stat mathbattle.Stat) string

	// Replies used in CmdSendServiceMessage
	ServiceMsgGetText() string
	ServiceMsgTextIsEmpty() string
	ServiceMsgCancelSend() string
	ServiceMsgAskRecieversType() string
	ServiceMsgWrongRecieversType() string
	ServiceMsgRecieversTypeAll() string
	ServiceMsgRecieversTypeSome() string
	ServiceMsgInputRecievers() string
	ServiceMsgFinalAsk(recieversType string, recievers ...string) string
	ServiceMsgSendSuccess() string

	// Replies used in CmdGetMyResults
	MyResultsProblemNotSolved(problemCaption string) string
	MyResultsProblemResults(problemCaption string, juriComment string, mark mathbattle.Mark,
		otherParticipantsReviews []mathbattle.Review) string
	MyResultsReviewResults(problemCaption string, solutionNumber int, isCommented bool, juriComment string, mark mathbattle.Mark) string
}
