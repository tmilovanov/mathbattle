package replier

import (
	"fmt"

	mathbattle "mathbattle/models"
)

type RussianReplier struct{}

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

func (r RussianReplier) Yes() string {
	return "Да"
}

func (r RussianReplier) No() string {
	return "Нет"
}

func (r RussianReplier) Cancel() string {
	return "Отменено"
}

func (r RussianReplier) GetStartMessage() string {
	return "Привет! Этот бот позволяет тебе участвовать в математических боях.\n"
}

func (r RussianReplier) GetHelpMessage(availableCommands []mathbattle.TelegramCommandHandler) string {
	msg := "Сейчас тебе доступны следующие действия:\n"
	msg += "\n"
	for _, cmd := range availableCommands {
		msg += cmd.Name() + " " + cmd.Description() + "\n"
	}
	return msg
}

func (r RussianReplier) CmdSubscribeName() string {
	return "/subscribe"
}

func (r RussianReplier) CmdSubscribeDesc() string {
	return "Подписаться на рассылку задач"
}

func (r RussianReplier) CmdUnsubscribeName() string {
	return "/unsubscribe"
}

func (r RussianReplier) CmdUnsubscribeDesc() string {
	return "Отписаться от рассылки задач"
}

func (r RussianReplier) CmdGetProblemsName() string {
	return "/get_problems"
}

func (r RussianReplier) CmdGetProblemsDesc() string {
	return "Показать задачи"
}

func (r RussianReplier) CmdSubmitSolutionName() string {
	return "/submit_solution"
}

func (r RussianReplier) CmdSubmitSolutionDesc() string {
	return "Отправить решение на задачу из текущего раунда"
}

func (r RussianReplier) CmdStartReviewStageName() string {
	return "/start_review_stage"
}

func (r RussianReplier) CmdStartReviewStageDesc() string {
	return "Начать этап ревью"
}

func (r RussianReplier) InternalError() string {
	return "Произошла внутрення ошибка. Свяжись с организаторами и опиши свою проблему."
}

func (r RussianReplier) AlreadyRegistered() string {
	return "Ты уже подписан на рассылку задач."
}

func (r RussianReplier) RegisterNameExpect() string {
	return "Введи своё имя. Имя должно состоять только из букв."
}

func (r RussianReplier) RegisterNameWrong() string {
	return "Имя должно состоять только из букв."
}

func (r RussianReplier) RegisterGradeExpect() string {
	return "Отлично. Теперь укажи класс, в котором ты учишься (В ответе используй только цифры)"
}

func (r RussianReplier) RegisterGradeWrong() string {
	return "Введён неправильный класс. Ожидается число от 1 до 11"
}

func (r RussianReplier) RegisterSuccess() string {
	return "Ты успешно зарегистрирован. Теперь ожидай рассылки задач"
}

func (r RussianReplier) NotSubscribed() string {
	return "Ты не подписан на рассылку задач."
}

func (r RussianReplier) UnsubscribeSuccess() string {
	return "Ты успешно отписан от рассылки задач."
}

func (r RussianReplier) ProblemsPostBefore() string {
	return "Привет! А вот и задачи"
}

func (r RussianReplier) ProblemsPostAfter() string {
	msg := "Как будешь готов - присылай решение. Для этого жми сюда: \n"
	msg += r.CmdSubmitSolutionName() + "\n"
	msg += "Решение необходимо присылать фотографиями. "
	return msg
}

func (r RussianReplier) SolutionPartUploaded(partNumber int) string {
	return fmt.Sprintf("Завершена загрузка листа №%d", partNumber)
}

func (r RussianReplier) SolutionUploadSuccess(totalUpload int) string {
	return fmt.Sprintf("Загрузка решения завершена. Всего в твоём решении %d %s", totalUpload, GetDeclensionByNumeral([3]string{
		"лист", "листа", "листов",
	}, totalUpload))
}

func (r RussianReplier) SolutionExpectPart() string {
	msg := "Отлично, теперь посылай решение. Решение необходимо отправлять фотографиями. "
	msg += "Ты можешь загрузить сколько угодно фотографий. "
	msg += "После того как отошлёшь всё - нажми кнопку '" + r.SolutionFinishUploading() + "'"
	return msg
}

func (r RussianReplier) SolutionIsRewriteOld() string {
	msg := "Для этой задачи ты уже отправлял решение. Новое решение перезапишет старое.\n"
	msg += "\n"
	msg += "Продолжить?"
	return msg
}

func (r RussianReplier) SolutionDeclineRewriteOld() string {
	return "Отменено"
}

func (r RussianReplier) SolutionWrongProblemNumberFormat() string {
	return "SolutionWrongProblemNumberFormat"
}

func (r RussianReplier) SolutionWrongProblemNumber() string {
	return "Неверный номер задачи."
}

func (r RussianReplier) SolutionExpectProblemNumber() string {
	return "Введи номер задачи, для которой ты хочешь отправить решение."
}

func (r RussianReplier) SolutionFinishUploading() string {
	return "Завершить отправку решения"
}

func (r RussianReplier) SolutionWrongFormat() string {
	return "Неверный формат решения. В решении ожидаются только фотографии"
}

func (r RussianReplier) SolutionEmpty() string {
	return "Ты не отправил ни одной фотографии своего решения :("
}

func (r RussianReplier) ReviewPost() string {
	return "Это решение другого участника, в котором ты должен отыскать недочёты."
}
