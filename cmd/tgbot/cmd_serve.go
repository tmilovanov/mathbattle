package main

import (
	"log"
	"strconv"
	"time"

	"mathbattle/cmd/tgbot/handlers"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	solutiondist "mathbattle/solution_distributor"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TelegramPostman struct {
	bot          *tb.Bot
	participants mathbattle.ParticipantRepository
}

func (pm *TelegramPostman) Post(participantID string, m *tb.Message) error {
	participant, err := pm.participants.GetByID(participantID)
	if err != nil {
		return err
	}

	chatID, err := strconv.ParseInt(participant.TelegramID, 10, 64)
	if err != nil {
		return err
	}

	_, err = pm.bot.Send(tb.ChatID(chatID), m)
	return err
}

func createCommands(storage mathbattle.Storage, replier mreplier.Replier, postman mathbattle.TelegramPostman) []mathbattle.TelegramCommandHandler {
	solutionDistributor := solutiondist.SolutionDistributor{}

	commandStart := &handlers.Start{
		Handler:     handlers.Handler{Name: "/start", Description: ""},
		Replier:     replier,
		AllCommands: []mathbattle.TelegramCommandHandler{},
	}
	result := []mathbattle.TelegramCommandHandler{
		&handlers.Subscribe{
			Handler: handlers.Handler{
				Name:        replier.CmdSubscribeName(),
				Description: replier.CmdSubscribeDesc(),
			},
			Replier:      replier,
			Participants: storage.Participants,
		},
		&handlers.Unsubscribe{
			Handler: handlers.Handler{
				Name:        replier.CmdUnsubscribeName(),
				Description: replier.CmdUnsubscribeDesc(),
			},
			Replier:      replier,
			Participants: storage.Participants,
			Rounds:       storage.Rounds,
		},
		&handlers.SubmitSolution{
			Handler: handlers.Handler{
				Name:        replier.CmdSubmitSolutionName(),
				Description: replier.CmdSubmitSolutionDesc(),
			},
			Replier:      replier,
			Participants: storage.Participants,
			Rounds:       storage.Rounds,
			Solutions:    storage.Solutions,
		},
		&handlers.StartReviewStage{
			Handler: handlers.Handler{
				Name:        replier.CmdStartReviewStageName(),
				Description: replier.CmdStartReviewStageDesc(),
			},
			Replier:             replier,
			Rounds:              storage.Rounds,
			Solutions:           storage.Solutions,
			SolutionDistributor: &solutionDistributor,
			ReviewersCount:      2,
			Postman:             postman,
		},
		commandStart,
	}
	commandStart.AllCommands = result

	return result
}

func commandServe(storage mathbattle.Storage, token string, ctxRepository mathbattle.TelegramContextRepository, replier mreplier.Replier) {
	b, err := tb.NewBot(tb.Settings{
		Token:       token,
		Poller:      &tb.LongPoller{Timeout: 10 * time.Second},
		Synchronous: true,
		//Verbose:     true,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	allCommands := createCommands(storage, replier, &TelegramPostman{bot: b})

	genericHandler := func(handler mathbattle.TelegramCommandHandler, m *tb.Message, startType mathbattle.CommandStep) {
		ctx, err := ctxRepository.GetByTelegramID(int64(m.Sender.ID))
		isSuitable, err := handler.IsCommandSuitable(ctx)
		if err != nil {
			b.Send(m.Sender, replier.InternalError())
			log.Printf("Failed to get user context: %v", err)
			return
		}

		if !isSuitable {
			b.Send(m.Sender, replier.GetHelpMessage(mathbattle.FilterCommandsToShow(allCommands, ctx)))
			return
		}

		if err != nil {
			b.Send(m.Sender, replier.InternalError())
			log.Printf("Failed to get user context: %v", err)
			return
		}
		defer func() {
			ctxRepository.Update(int64(m.Sender.ID), ctx)
		}()

		if startType == mathbattle.StepStart {
			ctx.CurrentStep = 0
		}
		if startType == mathbattle.StepNext {
			ctx.CurrentStep = ctx.CurrentStep + 1
		}
		ctx.CurrentCommand = handler.Name()
		newStep, response, err := handler.Handle(ctx, m)
		if err != nil {
			b.Send(m.Sender, replier.InternalError())
			log.Printf("Failed to handle command: %s : %v", handler.Name(), err)
		}
		if response.Text != "" {
			b.Send(m.Sender, response.Text, response.Keyboard)
			if newStep == -1 && err == nil { // Command finished
				b.Send(m.Sender, replier.GetHelpMessage(mathbattle.FilterCommandsToShow(allCommands, ctx)))
			}
		}

		if newStep == -1 {
			ctx.CurrentStep = 0
			ctx.CurrentCommand = ""
		} else {
			ctx.CurrentStep = newStep
		}
	}

	for _, handler := range allCommands {
		b.Handle(handler.Name(), func(handler mathbattle.TelegramCommandHandler) func(m *tb.Message) {
			return func(m *tb.Message) {
				genericHandler(handler, m, mathbattle.StepStart)
			}
		}(handler))
	}

	genericMessagesHandler := func(m *tb.Message) {
		ctx, err := ctxRepository.GetByTelegramID(int64(m.Sender.ID))
		if err != nil {
			b.Send(m.Sender, replier.InternalError())
			log.Printf("Failed to get user context: %v", err)
			return
		}

		for _, handler := range allCommands {
			if handler.Name() == ctx.CurrentCommand {

				fillFileStruct := func(f tb.File) (tb.File, error) {
					result := f

					tmp, err := b.FileByID(f.FileID)
					if err != nil {
						return f, err
					}
					result.FilePath = tmp.FilePath

					fileReader, err := b.GetFile(&f)
					if err != nil {
						return f, err
					}
					result.FileReader = fileReader

					return result, nil
				}

				if m.Photo != nil {
					m.Photo.File, err = fillFileStruct(m.Photo.File)
					if err != nil {
						b.Send(m.Sender, replier.InternalError())
						log.Printf("Failed to fill photo structure: %v", err)
					}
				}

				if m.Document != nil {
					m.Document.File, err = fillFileStruct(m.Document.File)
					if err != nil {
						b.Send(m.Sender, replier.InternalError())
						log.Printf("Failed to fill document structure: %v", err)
					}
				}

				genericHandler(handler, m, mathbattle.StepSame)

				return
			}
		}

		b.Send(m.Sender, replier.GetHelpMessage(mathbattle.FilterCommandsToShow(allCommands, ctx)))
	}

	b.Handle(tb.OnPhoto, genericMessagesHandler)
	b.Handle(tb.OnText, genericMessagesHandler)
	b.Handle(tb.OnDocument, genericMessagesHandler)

	log.Printf("Bot started")

	b.Start()
}
