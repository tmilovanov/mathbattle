package infrastructure

import (
	"log"

	"mathbattle/application"
	problemdistributor "mathbattle/application/problem_distributor"
	solutiondistributor "mathbattle/application/solution_distributor"
	"mathbattle/config"
	"mathbattle/infrastructure/repository/sqldb"
	"mathbattle/interfaces/replier"
	"mathbattle/models/mathbattle"
)

type Container struct {
	cfg config.Config

	// Server side services
	roundService       *application.RoundService
	statService        *application.StatService
	participantService *application.ParticipantService
	solutionService    *application.SolutionService
	reviewService      *application.ReviewService
	problemService     *application.ProblemService

	// Others
	replier                application.Replier
	userRepository         *sqldb.UserRepository
	participantRepsitory   *sqldb.ParticipantRepository
	roundRepository        *sqldb.RoundRepository
	problemRepository      *sqldb.ProblemRepository
	solutionRepository     *sqldb.SolutionRepository
	reviewRepository       *sqldb.ReviewRepository
	postman                mathbattle.Postman
	solveStageDistributor  application.ProblemDistributor
	reviewStageDistributor application.SolutionDistributor
}

func NewServerContainer(config config.Config) Container {
	return Container{
		cfg: config,
	}
}

func (c *Container) Config() config.Config {
	return c.cfg
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

func (c *Container) Replier() application.Replier {
	if c.replier == nil {
		c.replier = &replier.RussianReplier{}
	}

	return c.replier
}

func (c *Container) UserRepository() mathbattle.UserRepository {
	if c.userRepository == nil {
		var err error
		c.userRepository, err = sqldb.NewUserRepository(c.Config().DatabaseType, c.Config().DatabaseConnectionString)
		if err != nil {
			log.Fatalf("Failed to get user repository, error: %v", err)
		}

		c.participantRepsitory, err = sqldb.NewParticipantRepository(c.Config().DatabaseType, c.Config().DatabaseConnectionString, c.userRepository)
		if err != nil {
			log.Fatalf("Failed to get participant repository, error: %v", err)
		}

		c.userRepository.SetParticipantRepository(c.participantRepsitory)
	}

	return c.userRepository
}

func (c *Container) RoundRepository() mathbattle.RoundRepository {
	if c.roundRepository == nil {
		var err error
		c.roundRepository, err = sqldb.NewRoundRepository(c.Config().DatabaseType, c.Config().DatabaseConnectionString)
		if err != nil {
			log.Fatalf("Failed to get round repository, error: %v", err)
		}
	}

	return c.roundRepository
}

func (c *Container) ParticipantRepository() mathbattle.ParticipantRepository {
	if c.participantRepsitory == nil {
		if c.userRepository == nil {
			var err error
			c.userRepository, err = sqldb.NewUserRepository(c.Config().DatabaseType, c.Config().DatabaseConnectionString)
			if err != nil {
				log.Fatalf("Failed to initialize user repository, error: %v", err)
			}
		}

		var err error
		c.participantRepsitory, err = sqldb.NewParticipantRepository(c.Config().DatabaseType, c.Config().DatabaseConnectionString, c.userRepository)
		if err != nil {
			log.Fatalf("Failed to get participant repository, error: %v", err)
		}
	}

	return c.participantRepsitory
}

func (c *Container) ProblemRepository() mathbattle.ProblemRepository {
	if c.problemRepository == nil {
		var err error
		c.problemRepository, err = sqldb.NewProblemRepository(c.Config().DatabaseType, c.Config().DatabaseConnectionString, c.Config().ProblemsPath)
		if err != nil {
			log.Fatalf("Failed to get problems repository, error: %v", err)
		}
	}

	return c.problemRepository
}

func (c *Container) SolutionRepository() mathbattle.SolutionRepository {
	if c.solutionRepository == nil {
		var err error
		c.solutionRepository, err = sqldb.NewSolutionRepository(c.Config().DatabaseType, c.Config().DatabaseConnectionString, c.Config().SolutionsPath)
		if err != nil {
			log.Fatalf("Failed to get solutions repository, error: %v", err)
		}
	}

	return c.solutionRepository
}

func (c *Container) ReviewRepository() mathbattle.ReviewRepository {
	if c.reviewRepository == nil {
		var err error
		c.reviewRepository, err = sqldb.NewReviewRepository(c.Config().DatabaseType, c.Config().DatabaseConnectionString)
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