package replier

import (
	mathbattle "mathbattle/models"
)

type Replier interface {
	GetReply(replyType BotReply) string
	GetHelpMessage(availableCommands []mathbattle.TelegramCommandHandler) string
	GetReplySSolutionUploadSuccess(totalUpload int) string
	GetReplySSolutionPartUploaded(partNumber int) string
}

type BotReply string

const (
	ReplyInternalErrorHappened             BotReply = "replyInternalErrorHappened"
	ReplyParticipantHello                  BotReply = "replyHello"
	ReplyUnknownHello                      BotReply = "replyUnknownHello"
	ReplyAlreadyRegistered                 BotReply = "replyAlreadyRegistered"
	ReplyRegisterNameExpect                BotReply = "replyRegisterNameExpect"
	ReplyRegisterNameWrong                 BotReply = "replyRegisterNameWrong"
	ReplyRegisterGradeExpect               BotReply = "replyRegisterGradeExpect"
	ReplyRegisterGradeWrong                BotReply = "replyRegisterGradeWrong"
	ReplyRegisterSchoolExpect              BotReply = "replyRegisterSchoolExpect"
	ReplyRegisterSchoolWrong               BotReply = "replyRegisterSchoolWrong"
	ReplyRegisterGeneralError              BotReply = "replyRegisterGeneralError"
	ReplyRegisterSuccess                   BotReply = "replyRegisterSuccess"
	ReplyProblemsPost                      BotReply = "replyProblemsPost"
	ReplyWrongSolutionFormat               BotReply = "replyWrongSolutionFormat"
	ReplyUnsubscribeSuccess                BotReply = "replyUnsubscribeSuccess"
	ReplyUnsubscribeNotSubscribed          BotReply = "replyUnsubscribeNotSubscribed"
	ReplyYouAreNotRegistered               BotReply = "replyYouAreNotRegistered"
	ReplySSolutionExpectProblem            BotReply = "replySSolutionExpectProblem"
	ReplySSolutionWrongProblemNumberFormat BotReply = "replySSolutionWrongProblemNumberFormat"
	ReplySSolutionWrongProblemNumber       BotReply = "replySSolutionWrongProblemNumber"
	ReplySSolutionExpectStartAccept        BotReply = "replySSolutionExpectStartAccept"
	ReplySSoltuionFinishUploading          BotReply = "replySSoltuionFinishUploading"
	ReplySSoltuionPartUploadSuccess        BotReply = "replySSoltuionPartUploadSuccess"
	ReplySSoltuionUploadSuccess            BotReply = "replySSoltuionUploadSuccess"
)
