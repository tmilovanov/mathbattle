package main

import (
	"fmt"
	"log"
	"os"

	mreplier "mathbattle/cmd/tgbot/replier"
	"mathbattle/mocks"
	mathbattle "mathbattle/models"
	problemdist "mathbattle/problem_distributor"
	problemdistributor "mathbattle/problem_distributor"
	"mathbattle/repository/memory"
	"mathbattle/repository/sqlite"

	"gopkg.in/yaml.v2"
)

func main() {
	log.Printf("Application started, arguments: %v", os.Args)

	if len(os.Args) < 3 {
		fmt.Println("Expected command")
		os.Exit(1)
	}

	isDebug := false
	if os.Args[1] == "debug" {
		isDebug = true
	}

	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v\n", err)
	}

	telegramUserRepository, err := sqlite.NewTelegramUserRepository(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	storage := getStorage(cfg, isDebug)

	switch os.Args[2] {
	case "example":
		commandExample(cfg.Token)
	case "start-round":
		// Сейчас раунд добавляется "бесконечным". Добавить возможность передать срок окончания раунда
		commandStartRound(storage, cfg.Token, mreplier.RussianReplier{}, 2)
	case "delete-round":
		commandDeleteRound(storage)
	case "run":
		userCtxRepository, err := memory.NewTelegramContextRepository(&telegramUserRepository)
		if err != nil {
			log.Fatal(err)
		}
		problemDistributor := problemdist.NewSimpleDistributor(storage.Problems, 3)
		commandServe(storage, cfg.Token, &userCtxRepository, mreplier.RussianReplier{}, &problemDistributor)
	case "debug-run":
		userCtxRepository, err := memory.NewTelegramContextRepository(&telegramUserRepository)
		if err != nil {
			log.Fatal(err)
		}

		problemDistributor := problemdist.NewSimpleDistributor(storage.Problems, 3)
		mocks.GenReviewPendingRound(storage.Rounds, storage.Participants, storage.Solutions, storage.Problems,
			&problemdistributor.SimpleDistributor{}, 10, 3, []int{1, 3, 6})
		commandServe(storage, cfg.Token, &userCtxRepository, mreplier.RussianReplier{}, &problemDistributor)
	}
}

type config struct {
	Token         string `yaml:"token"`
	DatabasePath  string `yaml:"db_path"`
	ProblemsPath  string `yaml:"problems_path"`
	SolutionsPath string `yaml:"solutions_path"`
}

func getConfig() (config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return config{}, err
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}

func getStorage(cfg config, isDebug bool) mathbattle.Storage {
	storage := mathbattle.Storage{}
	if isDebug {
		dbName := "mathbattle_test.sqlite"
		problemsPath := "problems_test"
		solutionsPath := "solutions_test"

		participants, err := sqlite.NewParticipantRepositoryTemp(dbName)
		if err != nil {
			log.Fatal(err)
		}
		storage.Participants = &participants

		solutions, err := sqlite.NewSolutionRepositoryTemp(dbName, solutionsPath)
		if err != nil {
			log.Fatal(err)
		}
		storage.Solutions = &solutions

		problems, err := sqlite.NewProblemRepositoryTemp(dbName, problemsPath)
		if err != nil {
			log.Fatal(err)
		}
		storage.Problems = &problems

		rounds, err := sqlite.NewRoundRepositoryTemp(dbName)
		if err != nil {
			log.Fatal(err)
		}
		storage.Rounds = &rounds

		reviews, err := sqlite.NewReviewRepositoryTemp(cfg.DatabasePath)
		if err != nil {
			log.Fatal(err)
		}
		storage.Reviews = &reviews
	} else {
		participants, err := sqlite.NewParticipantRepository(cfg.DatabasePath)
		if err != nil {
			log.Fatal(err)
		}
		storage.Participants = &participants

		solutions, err := sqlite.NewSolutionRepository(cfg.DatabasePath, cfg.SolutionsPath)
		if err != nil {
			log.Fatal(err)
		}
		storage.Solutions = &solutions

		problems, err := sqlite.NewProblemRepository(cfg.DatabasePath, cfg.ProblemsPath)
		if err != nil {
			log.Fatal(err)
		}
		storage.Problems = &problems

		rounds, err := sqlite.NewRoundRepository(cfg.DatabasePath)
		if err != nil {
			log.Fatal(err)
		}
		storage.Rounds = &rounds

		reviews, err := sqlite.NewReviewRepository(cfg.DatabasePath)
		if err != nil {
			log.Fatal(err)
		}
		storage.Reviews = &reviews
	}

	return storage
}
