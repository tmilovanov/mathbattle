package main

// Ограничения Telegram для broadcasting
// https://core.telegram.org/bots/faq
// "The API will not allow bulk notifications to more than ~30 users per second"

import (
	"fmt"
	"log"
	"strconv"

	mreplier "mathbattle/cmd/tgbot/replier"
	"mathbattle/internal/distributor"
	mathbattle "mathbattle/models"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func commandStartRound(storage mathbattle.Storage, telegramToken string, replyer mreplier.Replier, problemCount int) {
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

	problemDistributor := distributor.RandomDistributor{}
	distribution, err := problemDistributor.Get(participants, problems, pastRounds, problemCount)
	if err != nil {
		log.Fatal(err)
	}

	round := mathbattle.NewRound()
	round.ProblemDistribution = distribution
	if err = storage.Rounds.Store(round); err != nil {
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

		if _, err = bot.Send(tgbotapi.NewMessage(chatID, replyer.GetReply(mreplier.ReplyProblemsPost))); err != nil {
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
	}
}
