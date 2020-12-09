package replier

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	mathbattle "mathbattle/models"
	"mathbattle/mstd"
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

func (r RussianReplier) GetAvailableCommands(availableCommands []mathbattle.TelegramCommandHandler) string {
	msg := "Сейчас доступны следующие действия:\n"
	msg += "\n"
	for _, cmd := range availableCommands {
		msg += cmd.Name() + " " + cmd.Description() + "\n"
	}
	return msg
}

func (r RussianReplier) GetHelpMessage() string {
	msg := `
	Этот бот поможет вам подготовиться к математическим боях. В процессе каждого раунда вы будете решать задачи, оформлять решения и проверять решения других участников. 
Чтобы принять участие, необходимо зарегистрироваться и подписаться на рассылку задач. Если не хотите участвовать в раунде, можете отписаться от рассылки задач.
Раунд состоит из трёх этапов: решение задач, проверка решений, подведение итогов. Каждый из этапов длится несколько дней, сроки объявляются в начале раунда. 

1. Решение задач. На старте раунда рассылаются задачи и объявляется срок, до которого следует сдать решения. Не обязательно решать все задачи, можете выбрать любые. Решение каждой задачи следует оформлять письменно на отдельном листе, затем отсканировать или качественно сфотографировать и отправить. 
Требования к решению:
•  Решение должно быть полным, все утверждения должны быть обоснованы. 
•  Решение должно быть аккуратно оформлено, разборчивым почерком.
•  Текст на фотографии должен быть чётким и легко читаться. Фотографируйте с хорошим освещением, перед отправкой отрегулируйте яркость и контрастность.
•  Фотографии должны быть правильно ориентированы. Если нужно, переверните перед отправкой. 
Чтобы оправить решение:
`
	msg += fmt.Sprintf("•  нажмите %s", r.CmdSubmitSolutionName())
	msg += `
•  выберете задачу.
•  загрузите изображения, их может быть несколько.
`
	msg += fmt.Sprintf("•  когда все изображения с решением данной задачи загружены, нажмите кнопку «%s»", r.SolutionFinishUploading())
	msg += `
•  после этого можно загрузить решение другой задачи. 
Если хотите поменять своё решение, отправьте новое на ту же задачу, новое решение заменит старое. После окончания этапа решения задач нельзя присылать и заменять решения.
2. Проверка решений. В начале этого этапа вы получите решения других участников. Это будут решения только тех задач, на которые вы сами отправили решение. Полученные решения необходимо проверить и написать комментарий. Если в одном решении несколько изображений, они приходят альбомом.
Требования к комментарию:
•  Начните с того, считаете вы решение верным или нет. Решение верное, если оно доведено до конца и все утверждения правильно обоснованы. 
•  Если считаете решение верным, можно ничего больше не писать или указать небольшие недочёты, если они есть
•  Если считаете решение неверным, объясните почему, укажите, где допущены ошибки. Можете предложить вариант исправления, но необязательно.
•  Не надо сравнивать решение со своим. Важна правильность проверяемого решения, а не его оптимальность или возможность других подходов. 
•  Запрещается критиковать автора решения и употреблять нецензурную лексику. 
Пример комментария. Я считаю решение неверным. В решении используется то, что 0,5n – целое число, но это верно только для чётных n, но по условию n – любое натуральное число. Для нечётных n утверждение не доказано.
Чтобы отправить комментарий:`
	msg += fmt.Sprintf("•  Нажмите %s", r.CmdSubmitReviewName())
	msg += `
•  Выберете номер решения, которое хотите прокомментировать
•  Напишите и отправьте комментарий
Комментарий можно изменить до окончания этапа проверки решений. Пожалуйста, старайтесь проверить все присланные решения. 
После окончания этапа комментарии отправляются авторам. Вам может прийти несколько комментариев на одно решение. 
3. Подведение итогов. Публикуются решения задач, вы можете изучить их и проверить себя. Присланные решение и комментарии просматриваются жюри, выставляются баллы. Жюри может отправить комментарий к решению. В тестовом раунде баллы отсутствуют, рейтинговая система будет доработана позже.
Если у вас есть вопросы, пишите @alzhukovskaia.
Спасибо за интерес!`
	return msg
}

func (r RussianReplier) CmdHelpName() string {
	return "/help"
}

func (r RussianReplier) CmdHelpDesc() string {
	return "Помощь"
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

func (r RussianReplier) SolutionWrongProblemCaptionFormat() string {
	return "Указана несуществующая задача."
}

func (r RussianReplier) SolutionWrongProblemCaption() string {
	return "Указана несуществующая задача."
}

func (r RussianReplier) SolutionExpectProblemCaption() string {
	return "Укажите задачу, для которой хотите отправить решение."
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
	hour, minute, sec := mstd.DurationToDayHourMinute(time.Until(untilDate))
	result += fmt.Sprintf("Общая продолжительность фазы отсылки ревью: %dд. %dч. %dм.\n", hour, minute, sec)
	result += "Верно?\n"
	return result
}

func (r RussianReplier) StartReviewSuccess() string {
	return "Решения разосланы, этап успешно начался."
}

func (r RussianReplier) ReviewPostBefore(stageDuration time.Duration, stageEnd time.Time) string {
	msg := "Начался этап взаимной проверки решений. "
	msg += "Во время него необходимо проверить решения других участников и найти в них недочёты, если они есть."

	day, hour, minute := mstd.DurationToDayHourMinute(stageDuration)
	msg += fmt.Sprintf("Этап продлится %dд. %dч. %dм. ", day, hour, minute)
	msg += fmt.Sprintf("После %s по московскому времени отправлять комментарии будет нельзя.", stageEnd.Format("02.01.2006 15:04"))

	msg += "\n"
	msg += "Ниже решения других участников, которые вам необходимо проверить."
	return msg
}

func (r RussianReplier) ReviewPostCaption(problemCaption string, solutionNumber int) string {
	return fmt.Sprintf("(Решение %d на задачу %s)", solutionNumber, problemCaption)
}

func (r RussianReplier) ReviewPostAfter() string {
	msg := "Как будете готовы - присылайте свои комментарии на решения других участников. Для этого нажмите сюда: \n"
	msg += r.CmdSubmitReviewName() + "\n"
	msg += "Комментарий следует присылать обычным текстом."
	return msg
}

func (r RussianReplier) ReviewGetSolutionCaptions(descriptors []mathbattle.SolutionDescriptor) []string {
	result := []string{}

	for _, descriptor := range descriptors {
		result = append(result, r.ReviewPostCaption(descriptor.ProblemCaption, descriptor.SolutionNumber))
	}

	return result
}

func (r RussianReplier) ReviewGetDescriptor(userInput string) (mathbattle.SolutionDescriptor, bool) {
	userInput = strings.Trim(userInput, "\t\r\n ")
	parts := strings.Split(userInput, " ")
	if len(parts) < 5 {
		return mathbattle.SolutionDescriptor{}, false
	}

	solutionNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return mathbattle.SolutionDescriptor{}, false
	}

	return mathbattle.SolutionDescriptor{
		ProblemCaption: parts[4][0 : len(parts[4])-1],
		SolutionNumber: solutionNumber,
	}, true
}

func (r RussianReplier) ReviewExpectSolutionCaption() string {
	return "Укажите решение, для которого хотите отправить комментарий."
}

func (r RussianReplier) ReviewWrongSolutionCaption() string {
	return "Указано несуществующее решение."
}

func (r RussianReplier) ReviewIsRewriteOld() string {
	msg := "Для этого решения вы уже отправляли комментарий. Новый комментарий перезапишет старый.\n"
	msg += "\n"
	msg += "Продолжить?"
	return msg
}

func (r RussianReplier) ReviewExpectContent() string {
	msg := "Отлично, теперь посылайте комментарий текстовым сообщением. "
	msg += "В комментарии следует указать, считаете ли вы решение верным. Если нет, объяснить почему."
	return msg
}

func (r RussianReplier) ReviewUploadSuccess() string {
	return "Комментарий записан."
}

func (r RussianReplier) ReviewMsgForReviewee(review mathbattle.Review) string {
	msg := "Вы получили комментарий на своё решение от другого участника:\n"
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
		days, hours, minutes := mstd.DurationToDayHourMinute(stat.TimeToSolveLeft)
		result += "Идёт фаза отсылки решений\n"
		result += fmt.Sprintf("До её конца осталось: %dд. %dч. %dм.\n", days, hours, minutes)
	} else if stat.RoundStage == mathbattle.StageReview {
		days, hours, minutes := mstd.DurationToDayHourMinute(stat.TimeToReviewLeft)
		result += "Идёт фаза проверки решений\n"
		result += fmt.Sprintf("До её конца осталось: %dд. %dч. %dм.\n", days, hours, minutes)
	}
	result += fmt.Sprintf("Всего решений прислано: %d\n", stat.SolutionsTotal)
	result += fmt.Sprintf("Всего комментариев прислано: %d\n", stat.ReviewsTotal)

	return result
}
