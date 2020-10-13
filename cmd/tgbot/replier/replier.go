package replier

import (
	mathbattle "mathbattle/models"
)

type Replier interface {
	GetHelpMessage(availableCommands []mathbattle.TelegramCommandHandler) string
	CmdSubscribeName() string
	CmdSubscribeDesc() string
	CmdUnsubscribeName() string
	CmdUnsubscribeDesc() string
	CmdSubmitSolutionName() string
	CmdSubmitSolutionDesc() string
	InternalError() string
	AlreadyRegistered() string
	RegisterNameExpect() string
	RegisterNameWrong() string
	RegisterGradeExpect() string
	RegisterGradeWrong() string
	RegisterSuccess() string
	NotSubscribed() string
	UnsubscribeSuccess() string
	ProblemsPost() string
	SolutionUploadSuccess(totalUpload int) string
	SolutionPartUploaded(partNumber int) string
	SolutionExpectProblemNumber() string
	SolutionWrongProblemNumberFormat() string
	SolutionWrongProblemNumber() string
	SolutionExpectPart() string
	SolutionFinishUploading() string
	SolutionWrongFormat() string
	SolutionEmpty() string
}
