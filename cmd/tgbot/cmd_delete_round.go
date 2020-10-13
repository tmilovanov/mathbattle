package main

import (
	"log"
	mathbattle "mathbattle/models"
)

func commandDeleteRound(storage mathbattle.Storage) {
	r, err := storage.Rounds.GetRunning()
	if err == mathbattle.ErrNotFound {
		return
	}

	if err != nil {
		log.Fatalf("Failed to get current round: %v", err)
	}

	if err := storage.Rounds.Delete(r.ID); err != nil {
		log.Fatalf("Failed to delete current round: %v", err)
	}

	solutions, err := storage.Solutions.FindMany(r.ID, "", "")
	if err != nil && err != mathbattle.ErrNotFound {
		log.Fatalf("Failed to find solutions for round: %v", err)
	}
	for _, solution := range solutions {
		err = storage.Solutions.Delete(solution.ID)
		if err != nil {
			log.Fatalf("Failed to delete solution for round: %v", err)
		}
	}
}
