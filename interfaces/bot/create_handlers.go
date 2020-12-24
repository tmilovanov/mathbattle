package bot

import (
	"mathbattle/infrastructure"
	"mathbattle/interfaces/bot/handlers"
)

func createCommands(container infrastructure.Container) []handlers.TelegramCommandHandler {
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
		&handlers.StartReviewStage{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdStartReviewStageName(),
				Description: container.Replier().CmdStartReviewStageDesc(),
			},
			Replier:      container.Replier(),
			RoundService: container.APIRoundService(),
		},
		&handlers.StartRound{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdStartRoundName(),
				Description: container.Replier().CmdStartRoundDesc(),
			},
			Replier:      container.Replier(),
			RoundService: container.APIRoundService(),
		},
		&handlers.Stat{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdStatName(),
				Description: container.Replier().CmdStatDesc(),
			},
			Replier:     container.Replier(),
			StatService: container.APIStatService(),
		},
		&handlers.Subscribe{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdSubscribeName(),
				Description: container.Replier().CmdSubscribeDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.APIParticipantService(),
		},
		&handlers.Unsubscribe{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdUnsubscribeName(),
				Description: container.Replier().CmdUnsubscribeDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.APIParticipantService(),
			RoundService:       container.APIRoundService(),
		},
		&handlers.SubmitSolution{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdSubmitSolutionName(),
				Description: container.Replier().CmdSubmitSolutionDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.APIParticipantService(),
			RoundService:       container.APIRoundService(),
			SolutionService:    container.APISolutionService(),
		},
		&handlers.SubmitReview{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdSubmitReviewName(),
				Description: container.Replier().CmdSubmitReviewDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.APIParticipantService(),
			RoundService:       container.APIRoundService(),
			ReviewService:      container.APIReviewService(),
		},
		&handlers.GetReviews{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdGetReviewsName(),
				Description: container.Replier().CmdGetReviewsDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.APIParticipantService(),
			ReviewService:      container.APIReviewService(),
			RoundService:       container.APIRoundService(),
			SolutionService:    container.APISolutionService(),
		},
		&handlers.GetProblems{
			Handler: handlers.Handler{
				Name:        container.Replier().CmdGetProblemsName(),
				Description: container.Replier().CmdGetProblemsDesc(),
			},
			Replier:            container.Replier(),
			ParticipantService: container.APIParticipantService(),
			RoundService:       container.APIRoundService(),
			ProblemService:     container.APIProblemService(),
		},
		commandStart,
	}

	return result
}
