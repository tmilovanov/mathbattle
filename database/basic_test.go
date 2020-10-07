package database

import (
	"bytes"
	"path/filepath"
	"testing"

	mathbattle "mathbattle/models"
)

var dbPath = "../mathbattle.sqlite"
var solutionPath = "../solution_store"

func isEqual(a mathbattle.Solution, b mathbattle.Solution) bool {
	if a.ParticipantID != b.ParticipantID {
		return false
	}
	if a.RoundID != b.RoundID {
		return false
	}
	if a.ProblemID != b.ProblemID {
		return false
	}

	if len(a.Parts) != len(b.Parts) {
		return false
	}

	for i := 0; i < len(a.Parts); i++ {
		if a.Parts[i].Extension != b.Parts[i].Extension {
			return false
		}
		if !bytes.Equal(a.Parts[i].Content, b.Parts[i].Content) {
			return false
		}
	}

	return true
}

func TestSQLSolutionRepository(t *testing.T) {
	p1, _ := filepath.Abs(dbPath)
	p2, _ := filepath.Abs(solutionPath)
	rep, err := NewSQLSolutionRepository(p1, p2)
	t.Logf("Path to database: %s, %s", p1, p2)
	if err != nil {
		t.Fail()
	}

	newEmptySolution := mathbattle.Solution{
		RoundID:       "1",
		ParticipantID: "1",
		ProblemID:     "1",
	}

	_, err = rep.Store(newEmptySolution)
	if err != nil {
		t.Logf("Failed to store empty solution: %v", err)
		t.Fail()
	}

	solution, err := rep.Find("1", "1", "1")
	if err != nil {
		t.Logf("Failed to get back empty solution: %v", err)
		t.Fail()
	}

	if !isEqual(solution, newEmptySolution) {
		t.Logf("Solutions are not equal")
		t.Fail()
	}

	err = rep.Delete(solution.ID)
	if err != nil {
		t.Fail()
	}

	newSolution := mathbattle.Solution{
		RoundID:       "1",
		ParticipantID: "1",
		ProblemID:     "1",
		Parts: []mathbattle.Image{
			mathbattle.Image{Extension: ".jpg", Content: []byte("123456")},
			mathbattle.Image{Extension: ".png", Content: []byte("654321")},
		},
	}

	_, err = rep.Store(newSolution)
	if err != nil {
		t.Logf("Failed to store solution: %v", err)
		t.Fail()
	}

	solution, err = rep.Find("1", "1", "1")
	if err != nil {
		t.Logf("Failed to find solution: %v", err)
		t.Fail()
	}

	if !isEqual(solution, newSolution) {
		t.Logf("Solutions aren't equal")
		t.Fail()
	}

	_, err = rep.Find("1", "1", "2")
	if err != mathbattle.ErrSolutionNotFound {
		t.Logf("Don't expect error except ErrNotFound")
		t.Fail()
	}

	newPart := mathbattle.Image{
		Extension: ".jpg",
		Content:   []byte("55555"),
	}
	newSolution.Parts = append(newSolution.Parts, newPart)

	err = rep.AppendPart(solution.ID, newPart)
	if err != nil {
		t.Logf("Failed to append part: %v", err)
		t.Fail()
	}

	solution, err = rep.Get(solution.ID)
	if err != nil {
		t.Fail()
	}

	if !isEqual(solution, newSolution) {
		t.Fail()
	}

	err = rep.Delete(solution.ID)
	if err != nil {
		t.Fail()
	}

	_, err = rep.Get(solution.ID)
	if err == nil {
		t.Fail()
	}
}
