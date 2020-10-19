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

	problems, err := storage.Problems.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	participants, err := storage.Participants.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	pastRounds, err := storage.Rounds.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	problemDistributor := problemdist.RandomDistributor{}
	distribution, err := problemDistributor.Get(participants, problems, pastRounds, problemCount)
	if err != nil {
		log.Fatal(err)
	}

	duration, _ := time.ParseDuration("48h")
	round := mathbattle.NewRound(duration)
	round.ProblemDistribution = distribution
	round, err = storage.Rounds.Store(round)
	if err != nil {
		log.Fatalf("Failed to save round: %v", err)
	}

	for participantID, problemsIDs := range round.ProblemDistribution {
		log.Printf("%s - %v", participantID, problemsIDs)

		participant, err := storage.Participants.GetByID(participantID)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = bot.Send(tgbotapi.NewMessage(participant.TelegramID, replier.ProblemsPostBefore())); err != nil {
			log.Fatalf("Failed to send problem to participant: %v", err)
		}

		for i, problemID := range problemsIDs {
			problem, err := storage.Problems.GetByID(problemID)
			if err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewPhotoUpload(participant.TelegramID, tgbotapi.FileBytes{Name: "", Bytes: problem.Content})
			msg.Caption = fmt.Sprintf("%d", i+1)
			if _, err := bot.Send(msg); err != nil {
				log.Fatalf("Failed to send problem to participant: %v", err)
			}
		}

		if _, err = bot.Send(tgbotapi.NewMessage(participant.TelegramID, replier.ProblemsPostAfter())); err != nil {
			log.Fatalf("Failed to send problem to participant: %v", err)
		}
	}
}
