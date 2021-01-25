package bot

import (
	"encoding/hex"
	"log"
	"time"

	"mathbattle/infrastructure"
	"mathbattle/infrastructure/repository/memory"
	"mathbattle/interfaces/bot/handlers"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TelegramContextRepository(container infrastructure.MBotContainer) *memory.TelegramContextRepository {
	telegramContextRepository, err := memory.NewTelegramContextRepository(container.UserRepository())
	if err != nil {
		log.Fatalf("Failed to get telegram context repository, error: %v", err)
	}

	return telegramContextRepository
}

func Start(container infrastructure.MBotContainer) {
	b, err := tb.NewBot(tb.Settings{
		Token:       container.Config().TelegramToken,
		Poller:      &tb.LongPoller{Timeout: 10 * time.Second},
		Synchronous: true,
		//Verbose:     true,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b2, err := tgbotapi.NewBotAPI(container.Config().TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	allCommands := createCommands(container)
	ctxRepository := TelegramContextRepository(container)

	commandHandler := func(handler handlers.TelegramCommandHandler, m *tb.Message, startType handlers.CommandStep) {
		ctx, err := ctxRepository.GetByUserData(infrastructure.TelegramUserData{
			ChatID:    int64(m.Sender.ID),
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			Username:  m.Sender.Username,
		})

		isSuitable, reason, err := handler.IsCommandSuitable(ctx)
		if err != nil {
			b.Send(m.Sender, container.Replier().InternalError(), &tb.ReplyMarkup{
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
				container.Replier().GetAvailableCommands(handlers.FilterCommandsToShow(allCommands, ctx)),
				&tb.ReplyMarkup{
					ReplyKeyboardRemove: true,
				})

			return
		}

		if !ctx.User.IsAdmin && handler.IsAdminOnly() {
			b.Send(m.Sender,
				container.Replier().GetAvailableCommands(handlers.FilterCommandsToShow(allCommands, ctx)),
				&tb.ReplyMarkup{
					ReplyKeyboardRemove: true,
				})
			return
		}

		if err != nil {
			b.Send(m.Sender, container.Replier().InternalError(), &tb.ReplyMarkup{
				ReplyKeyboardRemove: true,
			})
			log.Printf("Failed to get user context: %v", err)
			return
		}
		defer func() {
			ctxRepository.Update(int64(m.Sender.ID), ctx)
		}()

		if startType == handlers.StepStart {
			ctx.CurrentStep = 0
		}
		if startType == handlers.StepNext {
			ctx.CurrentStep = ctx.CurrentStep + 1
		}
		ctx.CurrentCommand = handler.Name()
		newStep, response, err := handler.Handle(ctx, m)
		if err != nil {
			b.Send(m.Sender, container.Replier().InternalError())
			log.Printf("Failed to handle command: %s : %v", handler.Name(), err)
		}
		if len(response) != 0 {
			for _, item := range response {
				if len(item.Img.Content) > 0 {
					// message with photo
					msg := tgbotapi.NewPhotoUpload(ctx.User.TelegramID, tgbotapi.FileBytes{Name: "", Bytes: item.Img.Content})
					msg.Caption = item.Text
					// msg.ReplyMarkup = item.Keyboard
					b2.Send(msg)
				} else {
					// text message only
					_, err = b.Send(m.Sender, item.Text, item.Keyboard, tb.ModeMarkdown)
					if err != nil {
						log.Printf("Failed to send message, error: %v", err)
					}
				}
			}

			if newStep == -1 && err == nil { // Command finished
				b.Send(m.Sender,
					container.Replier().GetAvailableCommands(handlers.FilterCommandsToShow(allCommands, ctx)),
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
		b.Handle(handler.Name(), func(handler handlers.TelegramCommandHandler) func(m *tb.Message) {
			return func(m *tb.Message) {
				commandHandler(handler, m, handlers.StepStart)
			}
		}(handler))
	}

	genericMessagesHandler := func(m *tb.Message) {
		hm, _ := hex.DecodeString("f09f9281f09f8fbce2808de29980efb88f")

		ctx, err := ctxRepository.GetByUserData(infrastructure.TelegramUserData{
			ChatID:    int64(m.Sender.ID),
			FirstName: string(hm),
			LastName:  string(hm),
			Username:  m.Sender.Username,
		})
		if err != nil {
			b.Send(m.Sender, container.Replier().InternalError())
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
						b.Send(m.Sender, container.Replier().InternalError())
						log.Printf("Failed to fill photo structure: %v", err)
					}
				}

				if m.Document != nil {
					m.Document.File, err = fillFileStruct(m.Document.File)
					if err != nil {
						b.Send(m.Sender, container.Replier().InternalError())
						log.Printf("Failed to fill document structure: %v", err)
					}
				}

				commandHandler(handler, m, handlers.StepSame)

				return
			}
		}

		b.Send(m.Sender,
			container.Replier().GetAvailableCommands(handlers.FilterCommandsToShow(allCommands, ctx)),
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
