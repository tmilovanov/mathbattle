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
	return "Привет! Этот бот поможет подготовиться к математическим боям.\n"
}

func (r RussianReplier) GetHelpMessage(availableCommands []mathbattle.TelegramCommandHandler) string {
	msg := "Сейчас доступны следующие действия:\n"
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
	return "Произошла внутрення ошибка. Свяжитесь с организаторами и опишите свою проблему."
}

func (r RussianReplier) AlreadyRegistered() string {
	return "Вы уже подписаны на рассылку задач."
}

func (r RussianReplier) RegisterNameExpect() string {
	return "Введите своё имя. Имя должно состоять только из букв."
}

func (r RussianReplier) RegisterNameWrong() string {
	return "Имя должно состоять только из букв."
}

func (r RussianReplier) RegisterGradeExpect() string {
	return "Отлично! Теперь укажите класс, в котором учитесь (используйте только цифры)."
}

func (r RussianReplier) RegisterGradeWrong() string {
	return "Введён неправильный класс. Ожидается число от 1 до 11."
}

func (r RussianReplier) RegisterSuccess() string {
	return "Вы успешно зарегистрированы. Ожидайте начала раунда и рассылки задач."
}

func (r RussianReplier) NotSubscribed() string {
	return "Вы не подписаны на рассылку задач."
}

func (r RussianReplier) UnsubscribeSuccess() string {
	return "Вы успешно отписаны от рассылки задач."
}

func (r RussianReplier) ProblemsPostBefore() string {
	return "Привет! А вот и задачи."
}

func (r RussianReplier) ProblemsPostAfter() string {
	msg := "Как будете готовы - присылайте решение. Для этого нажмите сюда: \n"
	msg += r.CmdSubmitSolutionName() + "\n"
	msg += "Решение следует оформить на бумаге и прислать качественное фото или скан. "
	return msg
}

func (r RussianReplier) SolutionPartUploaded(partNumber int) string {
	return fmt.Sprintf("Завершена загрузка листа №%d", partNumber)
}

func (r RussianReplier) SolutionUploadSuccess(totalUpload int) string {
	return fmt.Sprintf("Загрузка решения завершена. Всего в решении %d %s", totalUpload, GetDeclensionByNumeral([3]string{
		"лист", "листа", "листов",
	}, totalUpload))
}

func (r RussianReplier) SolutionExpectPart() string {
	msg := "Отлично, теперь присылайте решение. Решение следует оформить на бумаге и прислать качественное фото или скан."
	msg += "Можно загрузить сколько угодно фотографий. "
	msg += "После того как отошлёте всё - нажмите кнопку '" + r.SolutionFinishUploading() + "'"
	return msg
}

func (r RussianReplier) SolutionIsRewriteOld() string {
	msg := "Для этой задачи вы уже отправляли решение. Новое решение перезапишет старое.\n"
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
	return "Введите номер задачи, для которой хотите отправить решение."
}

func (r RussianReplier) SolutionFinishUploading() string {
	return "Завершить отправку решения"
}

func (r RussianReplier) SolutionWrongFormat() string {
	return "Неверный формат решения. В решении ожидаются только фотографии."
}

func (r RussianReplier) SolutionEmpty() string {
	return "Вы не отправили ни одной фотографии своего решения :("
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

func (r RussianReplier) ReviewPostBefore() string {
	return "Это решения другого участника, в котором нужно проверить и найти недочёты, если они есть."
}

func (r RussianReplier) ReviewPostCaption(problemIndex int, solutionNumber int) string {
	return fmt.Sprintf("%d (Решение на задачу %d)", solutionNumber, problemIndex)
}

func (r RussianReplier) ReviewPostAfter() string {
	msg := "Как будете готовы - присылайте свои отзывы на решения других участников. Для этого нажмите сюда: \n"
	msg += r.CmdSubmitReviewName() + "\n"
	msg += "Отзыв следует присылать обычным текстом."
	return msg
}

func (r RussianReplier) ReviewExpectSolutionNumber() string {
	return "Введите номер решения, для которого хотите отправить отзыв."
}

func (r RussianReplier) ReviewWrongSolutionNumber() string {
	return "Неверный номер решения."
}

func (r RussianReplier) ReviewIsRewriteOld() string {
	msg := "Для этой задачи вы уже отправляли отзыв. Новый отзыв перезапишет старый.\n"
	msg += "\n"
	msg += "Продолжить?"
	return msg
}

func (r RussianReplier) ReviewExpectContent() string {
	return "Отлично, теперь посылайте отзыв текстовым сообщением. В отзыве следует указать, считаете ли вы решение верным. Если нет, объяснить почему."
}

func (r RussianReplier) ReviewUploadSuccess() string {
	return "Отзыв записан."
}

func (r RussianReplier) ReviewMsgForReviewee(review mathbattle.Review) string {
	msg := "Вы получили отзыв на своё решение от другого участника:\n"
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
		result += "Идёт фаза проверки решений\n"
		result += fmt.Sprintf("До её конца осталось: %dд. %dч. %dм.\n", days, hours, minutes)
	}
	result += fmt.Sprintf("Всего решений прислано: %d\n", stat.SolutionsTotal)
	result += fmt.Sprintf("Всего ревью прислано: %d\n", stat.ReviewsTotal)

	return result
}
