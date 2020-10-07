package replier

import (
	mathbattle "mathbattle/models"
)

type RussianReplyer struct{}

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
