package replier

import (
	mathbattle "mathbattle/models"
	"mathbattle/usecases"
	"time"
)

type Replier interface {
	GetStartMessage() string
	GetHelpMessage(availableCommands []mathbattle.TelegramCommandHandler) string

	Yes() string
	No() string
	Cancel() string

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

	InternalError() string

	// Replies used in CmdSubscribe
	AlreadyRegistered() string
	RegisterNameExpect() string
	RegisterNameWrong() string
	RegisterGradeExpect() string
	RegisterGradeWrong() string
	RegisterSuccess() string

	// Replies used in CmdUnsubscribe
	NotSubscribed() string
	UnsubscribeSuccess() string

	// Replies used to post problems during start of round
	ProblemsPostBefore() string
	ProblemsPostAfter() string

	// Replies used in CmdSubmitSolution
	SolutionUploadSuccess(totalUpload int) string
	SolutionPartUploaded(partNumber int) string
	SolutionExpectProblemNumber() string
	SolutionWrongProblemNumberFormat() string
	SolutionWrongProblemNumber() string
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

	// Replies used to post solutions to other participants to review
	ReviewPostBefore() string
	ReviewPostCaption(problemIndex int, solutionNumber int) string
	ReviewPostAfter() string

	// Replies used in CmdSubmitReview
	ReviewExpectSolutionNumber() string
	ReviewWrongSolutionNumber() string
	ReviewIsRewriteOld() string
	ReviewExpectContent() string
	ReviewUploadSuccess() string
	ReviewMsgForReviewee(review mathbattle.Review) string

	// Replies used in CmdStat
	FormatStat(stat usecases.Stat) string
}
