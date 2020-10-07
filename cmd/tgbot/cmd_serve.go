package main

import (
	"log"
	"time"

	"mathbattle/cmd/tgbot/handlers"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

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

	commandStart := &handlers.Start{
		Handler:     handlers.Handler{Name: "/start", Description: ""},
		Replier:     replier,
		AllCommands: []mathbattle.TelegramCommandHandler{},
	}
	allCommands := []mathbattle.TelegramCommandHandler{
		&handlers.Subscribe{
			Handler:      handlers.Handler{Name: "/subscribe", Description: "Subscribe"},
			Replier:      replier,
			Participants: storage.Participants,
		},
		&handlers.Unsubscribe{
			Handler:      handlers.Handler{Name: "/unsubscribe", Description: "Unsubscribe"},
			Replier:      replier,
			Participants: storage.Participants,
			Rounds:       storage.Rounds,
		},
		&handlers.SubmitSolution{
			Handler:      handlers.Handler{Name: "/submit_solution", Description: "Submit solution"},
			Replier:      replier,
			Participants: storage.Participants,
			Rounds:       storage.Rounds,
			Solutions:    storage.Solutions,
		},
		commandStart,
	}
	commandStart.AllCommands = allCommands

	genericHandler := func(handler mathbattle.TelegramCommandHandler, m *tb.Message, startType mathbattle.CommandStep) {
		ctx, err := ctxRepository.GetByTelegramID(int64(m.Sender.ID), b)
		if err != nil {
			ctx.SendText(replier.GetReply(mreplier.ReplyInternalErrorHappened))
			log.Printf("Failed to get user context: %v", err)
			return
		}
		defer func() {
			ctxRepository.Update(ctx)
		}()

		if startType == mathbattle.StepStart {
			ctx.CurrentStep = 0
		}
		if startType == mathbattle.StepNext {
			ctx.CurrentStep = ctx.CurrentStep + 1
		}
		ctx.CurrentCommand = handler.Name()
		newStep, err := handler.Handle(ctx, m)
		if err != nil {
			if err == handlers.ErrCommandUnavailable {
				ctx.SendText(replier.GetHelpMessage(mathbattle.FilterCommandsToShow(allCommands, ctx)))
			} else {
				ctx.SendText(replier.GetReply(mreplier.ReplyInternalErrorHappened))
				log.Printf("Failed to handle command: %s : %v", handler.Name(), err)
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
		ctx, err := ctxRepository.GetByTelegramID(int64(m.Sender.ID), b)
		if err != nil {
			ctx.SendText(replier.GetReply(mreplier.ReplyInternalErrorHappened))
			log.Printf("Failed to get user context: %v", err)
			return
		}

		for _, handler := range allCommands {
			if handler.Name() == ctx.CurrentCommand {
				genericHandler(handler, m, mathbattle.StepSame)
				return
			}
		}

		ctx.SendText(replier.GetHelpMessage(mathbattle.FilterCommandsToShow(allCommands, ctx)))
	}

	b.Handle(tb.OnPhoto, genericMessagesHandler)
	b.Handle(tb.OnText, genericMessagesHandler)

	log.Printf("Bot started")

	b.Start()
}
