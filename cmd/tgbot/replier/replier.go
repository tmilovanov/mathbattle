package replier

import (
	mathbattle "mathbattle/models"
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
	CmdSubmitSolutionName() string
	CmdSubmitSolutionDesc() string
	CmdStartReviewStageName() string
	CmdStartReviewStageDesc() string

	InternalError() string
	AlreadyRegistered() string
	RegisterNameExpect() string
	RegisterNameWrong() string
	RegisterGradeExpect() string
	RegisterGradeWrong() string
	RegisterSuccess() string
	NotSubscribed() string
	UnsubscribeSuccess() string
	ProblemsPostBefore() string
	ProblemsPostAfter() string
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
	ReviewPost() string
}
