package replier

import (
	"fmt"
	mathbattle "mathbattle/models"
)

type RussianReplyer struct{}

func GetDeclensionByNumeral(forms [3]string, numeral int) string {
	ld := numeral % 10
	ltd := numeral % 100
	if ld == 1 && ltd != 11 {
		//1, 21, 31, ... 101, ...
		return forms[0]
	}

	if ld > 1 && ld < 5 && ltd != 12 && ltd != 13 && ltd != 14 {
		//2, 3, 4, 22, 23, 24, 32, 33, 34, ...
		return forms[1]
	}

	return forms[2]
}

func (r RussianReplyer) GetHelpMessage(availableCommands []mathbattle.TelegramCommandHandler) string {
	msg := "Привет! Этот бот позволяет тебе участвовать в математических боях. "
	msg += "Вот что ты можешь сейчас сделать:\n"
	msg += "\n"
	for _, cmd := range availableCommands {
		msg += cmd.Name() + " " + cmd.Description() + "\n"
	}
	return msg
}

func (r RussianReplyer) GetReply(replyType BotReply) string {
	switch replyType {
	case ReplyInternalErrorHappened:
		return "Произошла внутрення ошибка. Свяжись с организаторами и опиши свою проблему."
	case ReplyParticipantHello:
		return "Привет! Скоро будут рассылка задач."
	case ReplyUnknownHello:
		msg := "Мы тебя ещё не знаем. Если ты хочешь подписаться на рассылку задач, тебе нужно зарегистрироваться. "
		msg += "Для начала назови своё настоящее имя. (Имя должно состоять только из букв)"
		return msg
	case ReplyRegisterNameExpect:
		return "replyRegisterNameExpect"
	case ReplyRegisterNameWrong:
		return "Имя должно состоять только из букв"
	case ReplyRegisterGradeExpect:
		return "Отлично. Теперь укажи класс, в котором ты учишься (В ответе используй только цифры)"
	case ReplyRegisterGradeWrong:
		return "Введён неправильный класс. Ожидается число от 1 до 11"
	case ReplyRegisterSchoolExpect:
		return "replyRegisterSchoolExpect"
	case ReplyRegisterSchoolWrong:
		return "replyRegisterSchoolWrong"
	case ReplyRegisterGeneralError:
		return "Произошла какая-то странная ошибка. Обратись к организаторам напрямую"
	case ReplySSolutionExpectProblem:
		return "Введи номер задачи, для которой ты хочешь отправить решение. "
	case ReplySSolutionWrongProblemNumber:
		return "Неверный номер задачи. "
	case ReplySSolutionExpectStartAccept:
		msg := "Отлично, теперь посылай решение. Решение необходимо отправлять фотографиями. "
		msg += "Ты можешь загрузить сколько угодно фотографий. "
		msg += "После того как отошлёшь все нажми кнопку '" + r.GetReply(ReplySSoltuionFinishUploading) + "'"
		return msg
	case ReplySSoltuionFinishUploading:
		return "Завершить отправку решения"
	case ReplyRegisterSuccess:
		return "Ты успешно зарегистрирован. Теперь ожидай рассылки задач"
	case ReplyProblemsPost:
		msg := "Привет! А вот и задачи. Как будешь готов - присылай решение."
		msg += "Решение необходимо присылать фотографиями. "
		return msg
	case ReplyWrongSolutionFormat:
		msg := "Неверный формат решения. В решении ожидаются только фотографии"
		return msg
	default:
		return string(replyType)
	}
}

func (r RussianReplyer) GetReplySSolutionPartUploaded(partNumber int) string {
	return fmt.Sprintf("Завершена загрузка листа №%d", partNumber)
}

func (r RussianReplyer) GetReplySSolutionUploadSuccess(totalUpload int) string {
	return fmt.Sprintf("Загрузка решения завершена. Всего в твоём решении %d %s", totalUpload, GetDeclensionByNumeral([3]string{
		"лист", "листа", "листов",
	}, totalUpload))
}
