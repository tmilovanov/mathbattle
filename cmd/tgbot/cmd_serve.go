package main

import (
	"log"
	"time"

	"mathbattle/cmd/tgbot/handlers"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"mathbattle/repository/sqlite"
	"mathbattle/scheduler"
	solutiondist "mathbattle/solution_distributor"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
	tb "gopkg.in/tucnak/telebot.v2"
)

func createCommands(storage mathbattle.Storage, replier mreplier.Replier,
	postman mathbattle.TelegramPostman, problemDistributor mathbattle.ProblemDistributor,
	scheduler mathbattle.MessageScheduler) []mathbattle.TelegramCommandHandler {
	solutionDistributor := solutiondist.SolutionDistributor{}

	commandStart := &handlers.Start{
		Handler:     handlers.Handler{Name: "/start", Description: ""},
		Replier:     replier,
		AllCommands: []mathbattle.TelegramCommandHandler{},
	}
	result := []mathbattle.TelegramCommandHandler{
		&handlers.Help{
			Handler: handlers.Handler{
				Name:        replier.CmdHelpName(),
				Description: replier.CmdHelpDesc(),
			},
			Replier: replier,
		},
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
		&handlers.GetProblems{
			Handler: handlers.Handler{
				Name:        replier.CmdGetProblemsName(),
				Description: replier.CmdGetProblemsDesc(),
			},
			Replier:            replier,
			Participants:       storage.Participants,
			Rounds:             storage.Rounds,
			Problems:           storage.Problems,
			ProblemDistributor: problemDistributor,
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
			Participants:        storage.Participants,
			SolutionDistributor: &solutionDistributor,
			ReviewersCount:      2,
			Postman:             postman,
			Scheduler:           scheduler,
		},
		&handlers.SubmitReview{
			Handler: handlers.Handler{
				Name:        replier.CmdSubmitReviewName(),
				Description: replier.CmdSubmitReviewDesc(),
			},
			Replier:      replier,
			Participants: storage.Participants,
			Rounds:       storage.Rounds,
			Reviews:      storage.Reviews,
			Solutions:    storage.Solutions,
			Postman:      postman,
		},
		&handlers.Stat{
			Handler: handlers.Handler{
				Name:        replier.CmdStatName(),
				Description: replier.CmdStatDesc(),
			},
			Replier:      replier,
			Participants: storage.Participants,
			Rounds:       storage.Rounds,
			Solutions:    storage.Solutions,
			Reviews:      storage.Reviews,
		},
		&handlers.GetReviews{
			Handler: handlers.Handler{
				Name:        replier.CmdGetReviewsName(),
				Description: replier.CmdGetReviewsDesc(),
			},
			Replier:     replier,
			Participant: storage.Participants,
			Reviews:     storage.Reviews,
			Rounds:      storage.Rounds,
			Solutions:   storage.Solutions,
		},
		commandStart,
	}
	commandStart.AllCommands = result

	return result
}

func commandServe(storage mathbattle.Storage, databasePath string, token string, ctxRepository mathbattle.TelegramContextRepository,
	replier mreplier.Replier, problemDistributor mathbattle.ProblemDistributor) {
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

	b2, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
		return
	}

	postman := &TelegramPostman{bot: b}

	scheduledMessageRepository, err := sqlite.NewScheduledMessageRepository(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	scheduler := scheduler.NewMessageScheduler(&scheduledMessageRepository, storage.Participants, postman)
	err = scheduler.StartAll()
	if err != nil {
		log.Fatal(err)
	}

	allCommands := createCommands(storage, replier, postman, problemDistributor, &scheduler)

	genericHandler := func(handler mathbattle.TelegramCommandHandler, m *tb.Message, startType mathbattle.CommandStep) {
		ctx, err := ctxRepository.GetByTelegramID(int64(m.Sender.ID))
		isSuitable, reason, err := handler.IsCommandSuitable(ctx)
		if err != nil {
			b.Send(m.Sender, replier.InternalError(), &tb.ReplyMarkup{
				ReplyKeyboardRemove: true,
			})
			log.Printf("Failed to get user context: %v", err)
			return
		}

		if !isSuitable {
			if reason != "" {
				b.Send(m.Sender, reason, &tb.ReplyMarkup{
					ReplyKeyboardRemove: true,
				})
			}

			b.Send(m.Sender,
				replier.GetAvailableCommands(mathbattle.FilterCommandsToShow(allCommands, ctx)),
				&tb.ReplyMarkup{
					ReplyKeyboardRemove: true,
				})

			return
		}

		if !ctx.User.IsAdmin && handler.IsAdminOnly() {
			b.Send(m.Sender,
				replier.GetAvailableCommands(mathbattle.FilterCommandsToShow(allCommands, ctx)),
				&tb.ReplyMarkup{
					ReplyKeyboardRemove: true,
				})
			return
		}

		if err != nil {
			b.Send(m.Sender, replier.InternalError(), &tb.ReplyMarkup{
				ReplyKeyboardRemove: true,
			})
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
		if len(response) != 0 {
			for _, item := range response {
				if len(item.Img.Content) > 0 {
					// message with photo
					msg := tgbotapi.NewPhotoUpload(ctx.User.ChatID, tgbotapi.FileBytes{Name: "", Bytes: item.Img.Content})
					msg.Caption = item.Text
					// msg.ReplyMarkup = item.Keyboard
					b2.Send(msg)
				} else {
					// text message only
					b.Send(m.Sender, item.Text, item.Keyboard, tb.ModeMarkdown)
				}
			}

			if newStep == -1 && err == nil { // Command finished
				b.Send(m.Sender,
					replier.GetAvailableCommands(mathbattle.FilterCommandsToShow(allCommands, ctx)),
					&tb.ReplyMarkup{
						ReplyKeyboardRemove: true,
					})
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

		b.Send(m.Sender,
			replier.GetAvailableCommands(mathbattle.FilterCommandsToShow(allCommands, ctx)),
			&tb.ReplyMarkup{
				ReplyKeyboardRemove: true,
			})
	}

	b.Handle(tb.OnPhoto, genericMessagesHandler)
	b.Handle(tb.OnText, genericMessagesHandler)
	b.Handle(tb.OnDocument, genericMessagesHandler)

	log.Printf("Bot started")

	b.Start()
}
