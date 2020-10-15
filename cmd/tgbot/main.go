package main

import (
	"fmt"
	"log"
	"os"

	mreplier "mathbattle/cmd/tgbot/replier"
	"mathbattle/database/mem"
	"mathbattle/database/sqlite"
	mathbattle "mathbattle/models"

	"gopkg.in/yaml.v2"
)

func main() {
	log.Printf("Application started, arguments: %v", os.Args)

	if len(os.Args) < 2 {
		fmt.Println("Expected command")
		os.Exit(1)
	}

	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v\n", err)
	}

	participantRepository, err := sqlite.NewSQLParticipantRepository(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	solutionRepository, err := sqlite.NewSQLSolutionRepository(cfg.DatabasePath, cfg.SolutionsPath)
	if err != nil {
		log.Fatal(err)
	}
	problemRepository, err := sqlite.NewSQLProblemRepository(cfg.DatabasePath, cfg.ProblemsPath)
	if err != nil {
		log.Fatal(err)
	}
	roundRepository, err := sqlite.NewSQLRoundRepository(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	telegramUserRepository, err := sqlite.NewTelegramUserRepository(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	storage := mathbattle.Storage{
		Participants: &participantRepository,
		Rounds:       &roundRepository,
		Problems:     &problemRepository,
		Solutions:    &solutionRepository,
	}

	switch os.Args[1] {
	case "start-round":
		// Сейчас раунд добавляется "бесконечным". Добавить возможность передать срок окончания раунда
		commandStartRound(storage, cfg.Token, mreplier.RussianReplier{}, 2)
	case "delete-round":
		commandDeleteRound(storage)
	case "run":
		userCtxRepository, err := mem.NewTelegramContextRepository(&telegramUserRepository)
		if err != nil {
			log.Fatal(err)
		}
		commandServe(storage, cfg.Token, &userCtxRepository, mreplier.RussianReplier{})
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
