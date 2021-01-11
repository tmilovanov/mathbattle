package infrastructure

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"mathbattle/application"
	problemdistributor "mathbattle/application/problem_distributor"
	solutiondistributor "mathbattle/application/solution_distributor"
	"mathbattle/config"
	"mathbattle/infrastructure/repository/sqlite"
	"mathbattle/interfaces/replier"
	"mathbattle/libs/fstraverser"
	"mathbattle/libs/mstd"
	"mathbattle/mocks"
	"mathbattle/models/mathbattle"
)

type TestRoundDescription struct {
	ParticipantsCount int
	ProblemsOnEach    int
}

type TestContainer struct {
	cfg                config.TestConfig
	roundService       mathbattle.RoundService
	statService        mathbattle.StatService
	participantService mathbattle.ParticipantService
	solutionService    mathbattle.SolutionService

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

func NewTestContainer() TestContainer {
	log.Printf("NewTestContainer()")
	var err error

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current user: %v %v %v %v", usr.Name, usr.Username, usr.Uid, usr.Gid)

	testStoragePath := filepath.Join(usr.HomeDir, "mathbattle_test_storage")
	cfg := config.TestConfig{
		DatabasePath:  filepath.Join(testStoragePath, "test_mathbattle.sqlite"),
		ProblemsPath:  filepath.Join(testStoragePath, "test_problems"),
		SolutionsPath: filepath.Join(testStoragePath, "test_solutions"),
	}

	if _, err := os.Stat(testStoragePath); os.IsNotExist(err) {
		return TestContainer{
			cfg: cfg,
		}
	}

	info, err := os.Stat(testStoragePath)
	if err != nil {
		log.Fatalf("Failed to get info about database, error: %v", err)
	}

	log.Printf("Database permissions: %v", info.Mode().String())

	fstraverser.TraverseStartingFrom(testStoragePath, func(fi fstraverser.FileInformation) {
		//log.Printf("---> %s", fi.Path)
	})

	log.Printf("Removing %v", testStoragePath)
	err = sqlite.Deinit()
	if err != nil {
		log.Fatalf("Failed to deinit database, err: %v", err)
	}

	err = os.RemoveAll(testStoragePath)
	if err != nil {
		log.Fatalf("Failed to remove test storage path: %v, error: %v", testStoragePath, err)
	}

	return TestContainer{
		cfg: cfg,
	}
}

func (c *TestContainer) Config() config.TestConfig {
	return c.cfg
}

func (c *TestContainer) RoundService() mathbattle.RoundService {
	if c.roundService == nil {
		c.roundService = &application.RoundService{
			Rep:                    c.RoundRepository(),
			Replier:                c.Replier(),
			Postman:                c.Postman(),
			Participants:           c.ParticipantRepository(),
			Solutions:              c.SolutionRepository(),
			SolveStageDistributor:  c.SolveStageDistributor(),
			ReviewStageDistributor: c.ReviewStageDistributor(),
			ReviewersCount:         2,
		}
	}

	return c.roundService
}

func (c *TestContainer) StatService() mathbattle.StatService {
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

func (c *TestContainer) ParticipantService() mathbattle.ParticipantService {
	if c.participantService == nil {
		c.participantService = &application.ParticipantService{
			Rep: c.ParticipantRepository(),
		}
	}

	return c.participantService
}

func (c *TestContainer) SolutionService() mathbattle.SolutionService {
	if c.solutionService == nil {
		c.solutionService = &application.SolutionService{
			Rep:    c.SolutionRepository(),
			Rounds: c.RoundRepository(),
		}
	}

	return c.solutionService
}

func (c *TestContainer) Replier() application.Replier {
	if c.replier == nil {
		c.replier = &replier.RussianReplier{}
	}

	return c.replier
}

func (c *TestContainer) UserRepository() mathbattle.UserRepository {
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

func (c *TestContainer) RoundRepository() mathbattle.RoundRepository {
	if c.roundRepository == nil {
		var err error
		c.roundRepository, err = sqlite.NewRoundRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get round repository, error: %v", err)
		}
	}

	return c.roundRepository
}

func (c *TestContainer) ParticipantRepository() mathbattle.ParticipantRepository {
	if c.participantRepsitory == nil {
		if c.userRepository == nil {
			var err error
			c.userRepository, err = sqlite.NewUserRepository(c.Config().DatabasePath)
			if err != nil {
				log.Fatalf("TestContainer::ParticipantRepository(), failed to initialize user repository, error: %v", err)
			}
		}

		var err error
		c.participantRepsitory, err = sqlite.NewParticipantRepository(c.Config().DatabasePath, c.userRepository)
		if err != nil {
			log.Fatalf("TestContainer::ParticipantRepository(), failed to get participant repository, error: %v", err)
		}
	}

	return c.participantRepsitory
}

func (c *TestContainer) ProblemRepository() mathbattle.ProblemRepository {
	if c.problemRepository == nil {
		var err error
		c.problemRepository, err = sqlite.NewProblemRepository(c.Config().DatabasePath, c.Config().ProblemsPath)
		if err != nil {
			log.Fatalf("Failed to get problems repository, error: %v", err)
		}
	}

	return c.problemRepository
}

func (c *TestContainer) SolutionRepository() mathbattle.SolutionRepository {
	if c.solutionRepository == nil {
		var err error
		c.solutionRepository, err = sqlite.NewSolutionRepository(c.Config().DatabasePath, c.Config().SolutionsPath)
		if err != nil {
			log.Fatalf("Failed to get solutions repository, error: %v", err)
		}
	}

	return c.solutionRepository
}

func (c *TestContainer) ReviewRepository() mathbattle.ReviewRepository {
	if c.reviewRepository == nil {
		var err error
		c.reviewRepository, err = sqlite.NewReviewRepository(c.Config().DatabasePath)
		if err != nil {
			log.Fatalf("Failed to get review repository, error: %v", err)
		}
	}

	return c.reviewRepository
}

func (c *TestContainer) Postman() mathbattle.Postman {
	if c.postman == nil {
		c.postman = nil
	}

	return c.postman
}

func (c *TestContainer) SolveStageDistributor() application.ProblemDistributor {
	if c.solveStageDistributor == nil {
		result := problemdistributor.NewSimpleDistributor(c.ProblemRepository(), 3)
		c.solveStageDistributor = &result
	}

	return c.solveStageDistributor
}

func (c *TestContainer) ReviewStageDistributor() application.SolutionDistributor {
	if c.reviewStageDistributor == nil {
		result := solutiondistributor.SolutionDistributor{}
		c.reviewStageDistributor = &result
	}

	return c.reviewStageDistributor
}

func (c *TestContainer) CreateUsers(count int) {
	for i := 0; i < count; i++ {
		_, err := c.UserRepository().Store(mathbattle.User{
			TelegramID:       int64(i),
			TelegramName:     fmt.Sprintf("FakeTelegramUserName_%d", i),
			IsAdmin:          false,
			RegistrationTime: time.Now(),
		})
		if err != nil {
			log.Fatalf("TestContainer::CreateUsers failed, error: %v", err)
		}
	}
}

func (c *TestContainer) CreateParticipants(count int) {
	users, err := c.UserRepository().GetAll()
	if err != nil {
		log.Fatalf("TestContainer::CreateParticipants, failed to get all users, error: %v", err)
	}

	if len(users) < count {
		log.Fatalf("TestContainer::CreateParticipants, not enough users, to create participants")
	}

	for i := 0; i < count; i++ {
		_, err := c.ParticipantRepository().Store(mathbattle.Participant{
			User:     users[i],
			Name:     fmt.Sprintf("FakeName_%d", i),
			School:   "FakeSchool",
			Grade:    11,
			IsActive: true,
		})

		if err != nil {
			log.Fatalf("TestContainer::CreateParticipants, failed to store participant, err: %v", err)
		}
	}
}

func (c *TestContainer) CreateSolveStageRound(desc TestRoundDescription) mathbattle.Round {
	for _, problem := range mocks.GenProblems(desc.ProblemsOnEach, 1, 11) {
		_, err := c.ProblemRepository().Store(problem)
		if err != nil {
			log.Fatal(err)
		}
	}

	c.CreateUsers(desc.ParticipantsCount)
	c.CreateParticipants(desc.ParticipantsCount)

	allParticipants, err := c.ParticipantRepository().GetAll()
	if err != nil {
		log.Fatal(err)
	}

	round := mathbattle.Round{
		ProblemDistribution: make(map[string][]mathbattle.ProblemDescriptor),
	}
	round.SetSolveStartDate(time.Now())

	for _, participant := range allParticipants {
		problems, err := c.SolveStageDistributor().GetForParticipantCount(participant, desc.ProblemsOnEach)
		if err != nil {
			log.Fatal(err)
		}

		descriptors := []mathbattle.ProblemDescriptor{}
		for i, problem := range problems {
			descriptors = append(descriptors, mathbattle.ProblemDescriptor{
				Caption:   mstd.IndexToLetter(i),
				ProblemID: problem.ID,
			})
		}

		round.ProblemDistribution[participant.ID] = descriptors
	}

	round, err = c.RoundRepository().Store(round)
	if err != nil {
		log.Fatal(err)
	}

	return round
}
