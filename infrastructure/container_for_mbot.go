package infrastructure

import (
	"log"

	"mathbattle/application"
	problemdistributor "mathbattle/application/problem_distributor"
	solutiondistributor "mathbattle/application/solution_distributor"
	"mathbattle/config"
	"mathbattle/infrastructure/repository/sqlite"
	"mathbattle/interfaces/client"
	"mathbattle/interfaces/replier"
	"mathbattle/models/mathbattle"
)

type MBotContainer struct {
	cfg config.Config

	roundService       *client.APIRound
	statService        *client.APIStat
	participantService *client.APIParticipant
	solutionService    *client.APISolution
	reviewService      *client.APIReview
	problemService     *client.APIProblem

	replier                application.Replier
	userRepository         *sqlite.UserRepository
	participantRepsitory   *sqlite.ParticipantRepository
	roundRepository        *sqlite.RoundRepository
	problemRepository      *sqlite.ProblemRepository
	solutionRepository     *sqlite.SolutionRepository
	reviewRepository       *sqlite.ReviewRepository
	postman                mathbattle.Postman
	solveStageDistributor  application.ProblemDistributor
	reviewStageDistributor application.SolutionDistributor
}

func NewMBotContainer(config config.Config) MBotContainer {
	return MBotContainer{
		cfg: config,
	}
}

func (c *MBotContainer) Config() config.Config {
	return c.cfg
}

func (c *MBotContainer) APIBaseUrl() string {
	return "http://" + c.Config().APIUrl
}

func (c *MBotContainer) RoundService() mathbattle.RoundService {
	if c.roundService == nil {
		c.roundService = &client.APIRound{BaseUrl: c.APIBaseUrl()}
	}

	return c.roundService
}

func (c *MBotContainer) StatService() mathbattle.StatService {
	if c.statService == nil {
		c.statService = &client.APIStat{BaseUrl: c.APIBaseUrl()}
	}

	return c.statService
}

func (c *MBotContainer) ParticipantService() mathbattle.ParticipantService {
	if c.participantService == nil {
		c.participantService = &client.APIParticipant{BaseUrl: c.APIBaseUrl()}
	}

	return c.participantService
}

func (c *MBotContainer) SolutionService() mathbattle.SolutionService {
	if c.solutionService == nil {
		c.solutionService = &client.APISolution{BaseUrl: c.APIBaseUrl()}
	}

	return c.solutionService
}

func (c *MBotContainer) ReviewService() mathbattle.ReviewService {
	if c.reviewService == nil {
		c.reviewService = &client.APIReview{BaseUrl: c.APIBaseUrl()}
	}

	return c.reviewService
}

func (c *MBotContainer) ProblemService() mathbattle.ProblemService {
	if c.problemService == nil {
		c.problemService = &client.APIProblem{BaseUrl: c.APIBaseUrl()}
	}

	return c.problemService
}

func (c *MBotContainer) Replier() application.Replier {
	if c.replier == nil {
		c.replier = &replier.RussianReplier{}
	}

	return c.replier
}

func (c *MBotContainer) UserRepository() mathbattle.UserRepository {
	if c.userRepository == nil {
		var err error
		c.userRepository, err = sqlite.NewUserRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get user repository, error: %v", err)
		}

		c.participantRepsitory, err = sqlite.NewParticipantRepository(c.Config().DatabasePath, c.userRepository)
		if err != nil {
			log.Fatalf("Failed to get participant repository, error: %v", err)
		}

		c.userRepository.SetParticipantRepository(c.participantRepsitory)
	}

	return c.userRepository
}

func (c *MBotContainer) RoundRepository() mathbattle.RoundRepository {
	if c.roundRepository == nil {
		var err error
		c.roundRepository, err = sqlite.NewRoundRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get round repository, error: %v", err)
		}
	}

	return c.roundRepository
}

func (c *MBotContainer) ParticipantRepository() mathbattle.ParticipantRepository {
	if c.participantRepsitory == nil {
		if c.userRepository == nil {
			var err error
			c.userRepository, err = sqlite.NewUserRepository(c.Config().DatabasePath)
			if err != nil {
				log.Fatalf("Failed to initialize user repository, error: %v", err)
			}
		}

		var err error
		c.participantRepsitory, err = sqlite.NewParticipantRepository(c.Config().DatabasePath, c.userRepository)
		if err != nil {
			log.Fatalf("Failed to get participant repository, error: %v", err)
		}
	}

	return c.participantRepsitory
}

func (c *MBotContainer) ProblemRepository() mathbattle.ProblemRepository {
	if c.problemRepository == nil {
		var err error
		c.problemRepository, err = sqlite.NewProblemRepository(c.Config().DatabasePath, c.Config().ProblemsPath)
		if err != nil {
			log.Fatalf("Failed to get problems repository, error: %v", err)
		}
	}

	return c.problemRepository
}

func (c *MBotContainer) SolutionRepository() mathbattle.SolutionRepository {
	if c.solutionRepository == nil {
		var err error
		c.solutionRepository, err = sqlite.NewSolutionRepository(c.Config().DatabasePath, c.Config().SolutionsPath)
		if err != nil {
			log.Fatalf("Failed to get solutions repository, error: %v", err)
		}
	}

	return c.solutionRepository
}

func (c *MBotContainer) ReviewRepository() mathbattle.ReviewRepository {
	if c.reviewRepository == nil {
		var err error
		c.reviewRepository, err = sqlite.NewReviewRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get review repository, error: %v", err)
		}
	}

	return c.reviewRepository
}

func (c *MBotContainer) Postman() mathbattle.Postman {
	if c.postman == nil {
		result, err := NewTelegramPostman(c.Config().TelegramToken)
		if err != nil {
			log.Fatalf("Failed to create postman, error: %v", err)
		}
		c.postman = result
	}

	return c.postman
}

func (c *MBotContainer) SolveStageDistributor() application.ProblemDistributor {
	if c.solveStageDistributor == nil {
		result := problemdistributor.NewSimpleDistributor(c.ProblemRepository(), 3)
		c.solveStageDistributor = &result
	}

	return c.solveStageDistributor
}

func (c *MBotContainer) ReviewStageDistributor() application.SolutionDistributor {
	if c.reviewStageDistributor == nil {
		result := solutiondistributor.SolutionDistributor{}
		c.reviewStageDistributor = &result
	}

	return c.reviewStageDistributor
}
