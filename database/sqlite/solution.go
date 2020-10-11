package sqlite

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	mathbattle "mathbattle/models"
)

type SQLSolutionRepository struct {
	sqliteRepository
	solutionFolder string
}

func NewSQLSolutionRepository(dbPath, solutionPath string) (SQLSolutionRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return SQLSolutionRepository{}, err
	}

	return SQLSolutionRepository{
		sqliteRepository: sqliteRepository,
		solutionFolder:   solutionPath,
	}, nil
}

func (r *SQLSolutionRepository) getPartPath(solution mathbattle.Solution, i int, extension string) string {
	fileName := fmt.Sprintf("%s_%s_%s_%d%s", solution.RoundID, solution.ParticipantID,
		solution.ProblemID, i, extension)
	result := filepath.Join(r.solutionFolder, fileName)
	return result
}

func (r *SQLSolutionRepository) Store(solution mathbattle.Solution) (mathbattle.Solution, error) {
	result := solution

	extensions := ""
	if len(solution.Parts) > 0 {
		for i := 0; i < len(solution.Parts); i++ {
			partPath := r.getPartPath(solution, i, solution.Parts[i].Extension)
			err := ioutil.WriteFile(partPath, solution.Parts[i].Content, 0666)
			if err != nil {
				return result, err
			}
			extensions = extensions + solution.Parts[i].Extension + ","
		}
		extensions = extensions[:len(extensions)-1]
	}

	res, err := r.db.Exec("INSERT INTO solutions (round_id, participant_id, problem_id, parts) VALUES (?,?,?,?)",
		solution.RoundID, solution.ParticipantID, solution.ProblemID, extensions)
	if err != nil {
		return result, err
	}

	newID, err := res.LastInsertId()
	if err != nil {
		return result, err
	}
	result.ID = strconv.FormatInt(newID, 10)

	return result, nil
}

func (r *SQLSolutionRepository) Get(ID string) (mathbattle.Solution, error) {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return mathbattle.Solution{}, err
	}

	res := r.db.QueryRow("SELECT id, round_id, participant_id, problem_id, parts FROM solutions WHERE id = ?", intID)
	var partsExtensions string
	result := mathbattle.Solution{}

	err = res.Scan(&result.ID, &result.RoundID, &result.ParticipantID, &result.ProblemID, &partsExtensions)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

	if len(partsExtensions) == 0 {
		return result, nil
	}

	extensions := strings.Split(partsExtensions, ",")
	for i := 0; i < len(extensions); i++ {
		partPath := r.getPartPath(result, i, extensions[i])
		content, err := ioutil.ReadFile(partPath)
		if err != nil {
			return result, err
		}
		result.Parts = append(result.Parts, mathbattle.Image{
			Extension: extensions[i],
			Content:   content,
		})
	}

	return result, nil
}

func (r *SQLSolutionRepository) Find(roundID string, participantID string, problemID string) (mathbattle.Solution, error) {
	res := r.db.QueryRow("SELECT id FROM solutions WHERE round_id = ? AND participant_id = ? AND problem_id = ?",
		roundID, participantID, problemID)

	var intID int
	err := res.Scan(&intID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mathbattle.Solution{}, mathbattle.ErrNotFound
		}
		return mathbattle.Solution{}, err
	}

	ID := strconv.Itoa(intID)

	return r.Get(ID)
}

func (r *SQLSolutionRepository) FindOrCreate(roundID string, participantID string, problemID string) (mathbattle.Solution, error) {
	s, err := r.Find(roundID, participantID, problemID)
	if err == nil {
		return s, nil
	}

	if err != mathbattle.ErrNotFound {
		return mathbattle.Solution{}, err
	}

	return r.Store(mathbattle.Solution{
		ParticipantID: participantID,
		ProblemID:     problemID,
		RoundID:       roundID,
	})
}

func (r *SQLSolutionRepository) AppendPart(ID string, item mathbattle.Image) error {
	solution, err := r.Get(ID)
	if err != nil {
		return err
	}

	newPartID := len(solution.Parts)
	solution.Parts = append(solution.Parts, item)
	partPath := r.getPartPath(solution, newPartID, item.Extension)
	err = ioutil.WriteFile(partPath, item.Content, 0666)
	if err != nil {
		return err
	}

	extensions := ""
	for i := 0; i < len(solution.Parts); i++ {
		extensions = extensions + solution.Parts[i].Extension + ","
	}
	extensions = extensions[:len(extensions)-1]

	intID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("UPDATE solutions SET parts = ? WHERE id = ?", extensions, intID)

	return err
}

func (r *SQLSolutionRepository) Delete(ID string) error {
	solution, err := r.Get(ID)
	if err != nil {
		return err
	}

	for i := 0; i < len(solution.Parts); i++ {
		partPath := r.getPartPath(solution, i, solution.Parts[i].Extension)
		err = os.Remove(partPath)
		if err != nil {
			return err
		}
	}

	_, err = r.db.Exec("DELETE FROM solutions WHERE id = ?", solution.ID)
	return err
}
