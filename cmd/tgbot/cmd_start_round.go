package main

// Ограничения Telegram для broadcasting
// https://core.telegram.org/bots/faq
// "The API will not allow bulk notifications to more than ~30 users per second"

import (
	"log"
	"time"

	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"mathbattle/mstd"
	problemdist "mathbattle/problem_distributor"
	"mathbattle/repository/sqlite"
	"mathbattle/scheduler"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func commandStartRound(storage mathbattle.Storage, databasePath string, telegramToken string, replier mreplier.Replier, problemCount int) {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	participants, err := storage.Participants.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	problemDistributor := problemdist.NewSimpleDistributor(storage.Problems, problemCount)
	duration := time.Minute * 3
	round := mathbattle.NewRound(duration)
	endDate, err := mathbattle.ParseStageEndDate("18.12.2020")
	if err != nil {
		log.Fatal(err)
	}
	round.SetSolveEndDate(endDate)

	postman := &TelegramPostman2{bot: bot}

	scheduledMessageRepository, err := sqlite.NewScheduledMessageRepository(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	scheduler := scheduler.NewMessageScheduler(&scheduledMessageRepository, storage.Participants, postman)
	err = scheduler.Schedule(mathbattle.ScheduledMessage{
		Message:       replier.SolveStageEnd(),
		SendTime:      round.GetSolveEndDate(),
		RecieversType: mathbattle.Everyone,
		Recievers:     []string{},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, participant := range participants {
		participantProblems, err := problemDistributor.GetForParticipant(participant)
		if err != nil {
			log.Fatal(err)
		}

		for i, problem := range participantProblems {
			round.ProblemDistribution[participant.ID] = append(round.ProblemDistribution[participant.ID],
				mathbattle.ProblemDescriptor{
					Caption:   mstd.IndexToLetter(i),
					ProblemID: problem.ID,
				})
		}

		log.Printf("%s - %v", participant.ID, mathbattle.GetProblemIDs(participantProblems))

		duration := round.GetSolveStageDuration()
		stageEndMsk, err := round.GetSolveEndDateMsk()
		if err != nil {
			log.Fatal(err)
		}

		if _, err = bot.Send(tgbotapi.NewMessage(participant.TelegramID, replier.ProblemsPostBefore(duration, stageEndMsk))); err != nil {
			log.Fatalf("Failed to send problem to participant: %v", err)
		}

		for i := 0; i < len(participantProblems); i++ {
			msg := tgbotapi.NewPhotoUpload(participant.TelegramID, tgbotapi.FileBytes{Name: "", Bytes: participantProblems[i].Content})
			msg.Caption = round.ProblemDistribution[participant.ID][i].Caption

			if _, err := bot.Send(msg); err != nil {
				log.Fatalf("Failed to send problem to participant: %v", err)
			}
		}

		if _, err = bot.Send(tgbotapi.NewMessage(participant.TelegramID, replier.ProblemsPostAfter())); err != nil {
			log.Fatalf("Failed to send problem to participant: %v", err)
		}

	}

	round, err = storage.Rounds.Store(round)
	if err != nil {
		log.Fatalf("Failed to save round: %v", err)
	}

}
