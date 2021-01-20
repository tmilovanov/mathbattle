package bot

import (
	"mathbattle/infrastructure"
	"mathbattle/interfaces/bot/handlers"
)

func createCommands(container infrastructure.MBotContainer) []handlers.TelegramCommandHandler {
	commandStart := &handlers.Start{
		Handler: handlers.Handler{Name: "/start", Description: ""},
		Replier: container.Replier(),
	}

	result := []handlers.TelegramCommandHandler{
		&handlers.Help{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdHelpName(),
				Description: container.Replier().CmdHelpDesc(),
			},
			Replier: container.Replier(),
		},
		&handlers.SendServiceMessage{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdServiceMsgName(),
				Description: container.Replier().CmdGetMyResultsDesc(),
			},
			Replier:        container.Replier(),
			PostmanService: container.Postman(),
		},
		&handlers.StartReviewStage{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdStartReviewStageName(),
				Description: container.Replier().CmdStartReviewStageDesc(),
			},
			Replier:      container.Replier(),
			RoundService: container.RoundService(),
		},
		&handlers.StartRound{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdStartRoundName(),
				Description: container.Replier().CmdStartRoundDesc(),
			},
			Replier:      container.Replier(),
			RoundService: container.RoundService(),
		},
		&handlers.Stat{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdStatName(),
				Description: container.Replier().CmdStatDesc(),
			},
			Replier:     container.Replier(),
			StatService: container.StatService(),
		},
		&handlers.Subscribe{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdSubscribeName(),
				Description: container.Replier().CmdSubscribeDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.ParticipantService(),
		},
		&handlers.Unsubscribe{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdUnsubscribeName(),
				Description: container.Replier().CmdUnsubscribeDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.ParticipantService(),
			RoundService:       container.RoundService(),
		},
		&handlers.SubmitSolution{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdSubmitSolutionName(),
				Description: container.Replier().CmdSubmitSolutionDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.ParticipantService(),
			RoundService:       container.RoundService(),
			SolutionService:    container.SolutionService(),
		},
		&handlers.SubmitReview{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdSubmitReviewName(),
				Description: container.Replier().CmdSubmitReviewDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.ParticipantService(),
			RoundService:       container.RoundService(),
			ReviewService:      container.ReviewService(),
		},
		&handlers.GetReviews{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdGetReviewsName(),
				Description: container.Replier().CmdGetReviewsDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.ParticipantService(),
			ReviewService:      container.ReviewService(),
			RoundService:       container.RoundService(),
			SolutionService:    container.SolutionService(),
		},
		&handlers.GetProblems{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdGetProblemsName(),
				Description: container.Replier().CmdGetProblemsDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.ParticipantService(),
			RoundService:       container.RoundService(),
			ProblemService:     container.ProblemService(),
		},
		&handlers.GetMyResults{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdGetMyResultsName(),
				Description: container.Replier().CmdGetMyResultsDesc(),
			},
			Replier:            container.Replier(),
			RoundService:       container.RoundService(),
			SolutionService:    container.SolutionService(),
			ParticipantService: container.ParticipantService(),
			ReviewService:      container.ReviewService(),
		},
		commandStart,
	}

	return result
}
