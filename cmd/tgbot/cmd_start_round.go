package main

// Ограничения Telegram для broadcasting
// https://core.telegram.org/bots/faq
// "The API will not allow bulk notifications to more than ~30 users per second"

import (
	"fmt"
	"log"
	"time"

	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	problemdist "mathbattle/problem_distributor"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func commandStartRound(storage mathbattle.Storage, telegramToken string, replier mreplier.Replier, problemCount int) {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	participants, err := storage.Participants.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	problemDistributor := problemdist.NewSimpleDistributor(storage.Problems, 3)
	duration, _ := time.ParseDuration("48h")
	round := mathbattle.NewRound(duration)

	for _, participant := range participants {
		participantProblems, err := problemDistributor.GetForParticipant(participant)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s - %v", participant.ID, mathbattle.GetProblemIDs(participantProblems))

		if _, err = bot.Send(tgbotapi.NewMessage(participant.TelegramID, replier.ProblemsPostBefore())); err != nil {
			log.Fatalf("Failed to send problem to participant: %v", err)
		}

		for i, problem := range participantProblems {
			msg := tgbotapi.NewPhotoUpload(participant.TelegramID, tgbotapi.FileBytes{Name: "", Bytes: problem.Content})
			msg.Caption = fmt.Sprintf("%d", i+1)
			if _, err := bot.Send(msg); err != nil {
				log.Fatalf("Failed to send problem to participant: %v", err)
			}
		}

		if _, err = bot.Send(tgbotapi.NewMessage(participant.TelegramID, replier.ProblemsPostAfter())); err != nil {
			log.Fatalf("Failed to send problem to participant: %v", err)
		}

		round.ProblemDistribution[participant.ID] = mathbattle.GetProblemIDs(participantProblems)
	}

	round, err = storage.Rounds.Store(round)
	if err != nil {
		log.Fatalf("Failed to save round: %v", err)
	}

}
