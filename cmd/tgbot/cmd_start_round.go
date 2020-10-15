package main

// Ограничения Telegram для broadcasting
// https://core.telegram.org/bots/faq
// "The API will not allow bulk notifications to more than ~30 users per second"

import (
	"fmt"
	"log"
	"strconv"
	"time"

	mreplier "mathbattle/cmd/tgbot/replier"
	problemdist "mathbattle/internal/problem_distributor"
	mathbattle "mathbattle/models"

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

		chatID, err := strconv.ParseInt(participant.TelegramID, 10, 64)
		if err != nil {
			log.Fatalf("Failed to parse participant TelegramID: %v", err)
		}

		if _, err = bot.Send(tgbotapi.NewMessage(chatID, replier.ProblemsPostBefore())); err != nil {
			log.Fatalf("Failed to send problem to participant: %v", err)
		}

		for i, problemID := range problemsIDs {
			problem, err := storage.Problems.GetByID(problemID)
			if err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: "", Bytes: problem.Content})
			msg.Caption = fmt.Sprintf("%d", i+1)
			if _, err := bot.Send(msg); err != nil {
				log.Fatalf("Failed to send problem to participant: %v", err)
			}
		}

		if _, err = bot.Send(tgbotapi.NewMessage(chatID, replier.ProblemsPostAfter())); err != nil {
			log.Fatalf("Failed to send problem to participant: %v", err)
		}
	}
}
