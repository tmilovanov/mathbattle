package replier

import (
	"fmt"
	"time"

	mathbattle "mathbattle/models"
	"mathbattle/usecases"
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

func (r RussianReplier) CmdSubmitReviewName() string {
	return "/submit_review"
}

func (r RussianReplier) CmdSubmitReviewDesc() string {
	return "Отправить замечания по решению"
}

func (r RussianReplier) CmdStatName() string {
	return "/stat"
}

func (r RussianReplier) CmdStatDesc() string {
	return "Статистика"
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
	return "Неверный номер задачи."
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

func (r RussianReplier) StartReviewGetDuration() string {
	result := "Введите дату окончания раунда по московскому времени, в одном из следующих форматов:\n"
	result += "DD.MM.YYYY HH:MM (Ревью нельзя будет отослать после указанной даты)\n"
	result += "DD.MM.YYYY (Последний день приёма ревью. Приём ревью окончится в полночь)\n"
	return result
}

func (r RussianReplier) StartReviewWrongDuration() string {
	return "Дата окончания раунда введена неверно"
}

func (r RussianReplier) StartReviewConfirmDuration(untilDate time.Time) string {
	untillDateStr := untilDate.Format("02.01.2006 15:04")
	result := fmt.Sprintf("После %s ревью приниматься не будут\n", untillDateStr)
	hour, minute, sec := usecases.DurationToDayHourMinute(time.Until(untilDate))
	result += fmt.Sprintf("Общая продолжительность фазы отсылки ревью: %dд. %dч. %dм.\n", hour, minute, sec)
	result += "Верно?\n"
	return result
}

func (r RussianReplier) StartReviewSuccess() string {
	return "Решения разосланы, этап успешно начался."
}

func (r RussianReplier) ReviewPost() string {
	return "Это решения другого участника, в котором ты должен отыскать недочёты."
}

func (r RussianReplier) ReviewExpectSolutionNumber() string {
	return "Введи номер решения, для которого ты хочешь отправить отзыв."
}

func (r RussianReplier) ReviewWrongSolutionNumber() string {
	return "Неверный номер решения."
}

func (r RussianReplier) ReviewIsRewriteOld() string {
	msg := "Для этой задачи ты уже отправлял отзыв. Новый отзыв перезапишет старый.\n"
	msg += "\n"
	msg += "Продолжить?"
	return msg
}

func (r RussianReplier) ReviewExpectContent() string {
	return "Отлично, теперь посылай отзыв текстовым сообщением."
}

func (r RussianReplier) ReviewUploadSuccess() string {
	return "Отзыв записан."
}

func (r RussianReplier) ReviewMsgForReviewee(review mathbattle.Review) string {
	msg := "Ты получил критику своего решения от другого участника:\n"
	msg += "\n"
	msg += review.Content
	return msg
}

func (r RussianReplier) FormatStat(stat usecases.Stat) string {
	result := ""
	result += fmt.Sprintf("Участников всего: %d\n", stat.ParticipantsTotal)
	result += fmt.Sprintf("Из них новых сегодня: %d\n", stat.ParticipantsToday)

	if stat.RoundStage == mathbattle.StageNotStarted || stat.RoundStage == mathbattle.StageFinished {
		result += "Нет активного раунда."
		return result
	}

	result += "\nСтатистика активного раунда:\n"
	if stat.RoundStage == mathbattle.StageSolve {
		days, hours, minutes := usecases.DurationToDayHourMinute(stat.TimeToSolveLeft)
		result += "Идёт фаза отсылки решений\n"
		result += fmt.Sprintf("До её конца осталось: %dд. %dч. %dм.\n", days, hours, minutes)
	} else if stat.RoundStage == mathbattle.StageReview {
		days, hours, minutes := usecases.DurationToDayHourMinute(stat.TimeToReviewLeft)
		result += "Идёт фаза критики решений\n"
		result += fmt.Sprintf("До её конца осталось: %dд. %dч. %dм.\n", days, hours, minutes)
	}
	result += fmt.Sprintf("Всего решений прислано: %d\n", stat.SolutionsTotal)
	result += fmt.Sprintf("Всего ревью прилано: %d\n", stat.ReviewsTotal)

	return result
}
