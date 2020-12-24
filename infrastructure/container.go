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

type Container struct {
	cfg config.Config
	// Server side services
	roundService       mathbattle.RoundService
	statService        mathbattle.StatService
	participantService mathbattle.ParticipantService
	solutionService    mathbattle.SolutionService
	reviewService      mathbattle.ReviewService
	problemService     mathbattle.ProblemService

	// API services
	apiRoundService       mathbattle.RoundService
	apiStatService        mathbattle.StatService
	apiParticipantService mathbattle.ParticipantService
	apiSolutionService    mathbattle.SolutionService
	apiReviewService      mathbattle.ReviewService
	apiProblemService     mathbattle.ProblemService

	// Others
	replier                application.Replier
	userRepository         mathbattle.UserRepository
	participantRepsitory   mathbattle.ParticipantRepository
	roundRepository        mathbattle.RoundRepository
	problemRepository      mathbattle.ProblemRepository
	solutionRepository     mathbattle.SolutionRepository
	reviewRepository       mathbattle.ReviewRepository
	postman                mathbattle.Postman
	solveStageDistributor  application.ProblemDistributor
	reviewStageDistributor application.SolutionDistributor
}

func NewContainer(config config.Config) Container {
	return Container{
		cfg: config,
	}
}

func (c *Container) Config() config.Config {
	return c.cfg
}

func (c *Container) APIBaseUrl() string {
	return "http://" + c.Config().APIUrl
}

func (c *Container) RoundService() mathbattle.RoundService {
	if c.roundService == nil {
		result := &application.RoundService{
			Rep:                    c.RoundRepository(),
			Replier:                c.Replier(),
			Postman:                c.Postman(),
			Participants:           c.ParticipantRepository(),
			Solutions:              c.SolutionRepository(),
			SolveStageDistributor:  c.SolveStageDistributor(),
			ReviewStageDistributor: c.ReviewStageDistributor(),
			ReviewersCount:         2,
		}
		if err := result.StartSchedulingActions(); err != nil {
			log.Fatal(err)
		}
		c.roundService = result
	}

	return c.roundService
}

func (c *Container) StatService() mathbattle.StatService {
	if c.statService == nil {
		c.statService = &application.StatService{
			Participants: c.ParticipantRepository(),
			Rounds:       c.RoundRepository(),
			Solutions:    c.SolutionRepository(),
			Reviews:      c.ReviewRepository(),
		}
	}

	return c.statService
}

func (c *Container) ParticipantService() mathbattle.ParticipantService {
	if c.participantService == nil {
		c.participantService = &application.ParticipantService{
			Rep: c.ParticipantRepository(),
		}
	}

	return c.participantService
}

func (c *Container) SolutionService() mathbattle.SolutionService {
	if c.solutionService == nil {
		c.solutionService = &application.SolutionService{
			Rep:    c.SolutionRepository(),
			Rounds: c.RoundRepository(),
		}
	}

	return c.solutionService
}

func (c *Container) ReviewService() mathbattle.ReviewService {
	if c.reviewService == nil {
		c.reviewService = &application.ReviewService{
			Rep:       c.ReviewRepository(),
			Rounds:    c.RoundRepository(),
			Solutions: c.SolutionRepository(),
		}
	}

	return c.reviewService
}

func (c *Container) ProblemService() mathbattle.ProblemService {
	if c.problemService == nil {
		c.problemService = &application.ProblemService{
			Rep: c.ProblemRepository(),
		}
	}

	return c.problemService
}

func (c *Container) APIRoundService() mathbattle.RoundService {
	if c.apiRoundService == nil {
		c.apiRoundService = &client.APIRound{BaseUrl: c.APIBaseUrl()}
	}

	return c.apiRoundService
}

func (c *Container) APIStatService() mathbattle.StatService {
	if c.apiStatService == nil {
		c.apiStatService = &client.APIStat{BaseUrl: c.APIBaseUrl()}
	}

	return c.apiStatService
}

func (c *Container) APIParticipantService() mathbattle.ParticipantService {
	if c.apiParticipantService == nil {
		c.apiParticipantService = &client.APIParticipant{BaseUrl: c.APIBaseUrl()}
	}

	return c.apiParticipantService
}

func (c *Container) APISolutionService() mathbattle.SolutionService {
	if c.apiSolutionService == nil {
		c.apiSolutionService = &client.APISolution{BaseUrl: c.APIBaseUrl()}
	}

	return c.apiSolutionService
}

func (c *Container) APIReviewService() mathbattle.ReviewService {
	if c.apiReviewService == nil {
		c.apiReviewService = &client.APIReview{BaseUrl: c.APIBaseUrl()}
	}

	return c.apiReviewService
}

func (c *Container) APIProblemService() mathbattle.ProblemService {
	if c.apiProblemService == nil {
		c.apiProblemService = &client.APIProblem{BaseUrl: c.APIBaseUrl()}
	}

	return c.apiProblemService
}

func (c *Container) Replier() application.Replier {
	if c.replier == nil {
		c.replier = &replier.RussianReplier{}
	}

	return c.replier
}

func (c *Container) UserRepository() mathbattle.UserRepository {
	if c.userRepository == nil {
		var err error
		c.userRepository, err = sqlite.NewUserRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get user repository, error: %v", err)
		}
	}

	return c.userRepository
}

func (c *Container) RoundRepository() mathbattle.RoundRepository {
	if c.roundRepository == nil {
		var err error
		c.roundRepository, err = sqlite.NewRoundRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get round repository, error: %v", err)
		}
	}

	return c.roundRepository
}

func (c *Container) ParticipantRepository() mathbattle.ParticipantRepository {
	if c.participantRepsitory == nil {
		var err error
		c.participantRepsitory, err = sqlite.NewParticipantRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get participant repository, error: %v", err)
		}
	}

	return c.participantRepsitory
}

func (c *Container) ProblemRepository() mathbattle.ProblemRepository {
	if c.problemRepository == nil {
		var err error
		c.problemRepository, err = sqlite.NewProblemRepository(c.Config().DatabasePath, c.Config().ProblemsPath)
		if err != nil {
			log.Fatalf("Failed to get problems repository, error: %v", err)
		}
	}

	return c.problemRepository
}

func (c *Container) SolutionRepository() mathbattle.SolutionRepository {
	if c.solutionRepository == nil {
		var err error
		c.solutionRepository, err = sqlite.NewSolutionRepository(c.Config().DatabasePath, c.Config().SolutionsPath)
		if err != nil {
			log.Fatalf("Failed to get solutions repository, error: %v", err)
		}
	}

	return c.solutionRepository
}

func (c *Container) ReviewRepository() mathbattle.ReviewRepository {
	if c.reviewRepository == nil {
		var err error
		c.reviewRepository, err = sqlite.NewReviewRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get review repository, error: %v", err)
		}
	}

	return c.reviewRepository
}

func (c *Container) Postman() mathbattle.Postman {
	if c.postman == nil {
		result, err := NewTelegramPostman(c.Config().TelegramToken)
		if err != nil {
			log.Fatalf("Failed to create postman, error: %v", err)
		}
		c.postman = result
	}

	return c.postman
}

func (c *Container) SolveStageDistributor() application.ProblemDistributor {
	if c.solveStageDistributor == nil {
		result := problemdistributor.NewSimpleDistributor(c.ProblemRepository(), 3)
		c.solveStageDistributor = &result
	}

	return c.solveStageDistributor
}

func (c *Container) ReviewStageDistributor() application.SolutionDistributor {
	if c.reviewStageDistributor == nil {
		result := solutiondistributor.SolutionDistributor{}
		c.reviewStageDistributor = &result
	}

	return c.reviewStageDistributor
}
