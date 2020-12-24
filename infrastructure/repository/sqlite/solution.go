package sqlite

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"mathbattle/models/mathbattle"
)

type SolutionRepository struct {
	sqliteRepository
	solutionFolder string
}

func NewSolutionRepository(dbPath, solutionPath string) (*SolutionRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		if err := os.Mkdir(solutionPath, 0777); err != nil {
			return nil, err
		}
	}

	result := &SolutionRepository{
		sqliteRepository: sqliteRepository,
		solutionFolder:   solutionPath,
	}

	if err := result.CreateTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *SolutionRepository) CreateTable() error {
	createStmt := `CREATE TABLE IF NOT EXISTS solutions (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			round_id INTEGER,
			participant_id INTEGER,
			problem_id INTEGER,
			parts TEXT
		)`
	_, err := r.db.Exec(createStmt)
	return err
}

func (r *SolutionRepository) getPartPath(solution mathbattle.Solution, i int, extension string) string {
	fileName := fmt.Sprintf("%s_%s_%s_%d%s", solution.RoundID, solution.ParticipantID,
		solution.ProblemID, i, extension)
	result := filepath.Join(r.solutionFolder, fileName)
	return result
}

func (r *SolutionRepository) Store(solution mathbattle.Solution) (mathbattle.Solution, error) {
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

func (r *SolutionRepository) Get(ID string) (mathbattle.Solution, error) {
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

func (r *SolutionRepository) Find(roundID string, participantID string, problemID string) (mathbattle.Solution, error) {
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

func (r *SolutionRepository) FindMany(roundID string, participantID string, problemID string) ([]mathbattle.Solution, error) {
	query := "SELECT id, round_id, participant_id, problem_id, parts FROM solutions"
	whereClauses := []string{}
	whereArgs := []interface{}{}
	if roundID != "" {
		whereClauses = append(whereClauses, " round_id = ?")
		whereArgs = append(whereArgs, roundID)
	}
	if participantID != "" {
		whereClauses = append(whereClauses, " participant_id = ?")
		whereArgs = append(whereArgs, participantID)
	}
	if problemID != "" {
		whereClauses = append(whereClauses, " problem_id = ?")
		whereArgs = append(whereArgs, problemID)
	}
	if len(whereClauses) != 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := r.db.Query(query, whereArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []mathbattle.Solution{}, mathbattle.ErrNotFound
		}
	}
	defer rows.Close()

	result := []mathbattle.Solution{}
	for rows.Next() {
		curSolution := mathbattle.Solution{}
		var partsExtensions string
		err = rows.Scan(&curSolution.ID, &curSolution.RoundID, &curSolution.ParticipantID,
			&curSolution.ProblemID, &partsExtensions)
		if err != nil {
			if err == sql.ErrNoRows {
				return result, mathbattle.ErrNotFound
			}
			return result, err
		}

		if len(partsExtensions) != 0 {
			extensions := strings.Split(partsExtensions, ",")
			for i := 0; i < len(extensions); i++ {
				partPath := r.getPartPath(curSolution, i, extensions[i])
				content, err := ioutil.ReadFile(partPath)
				if err != nil {
					return result, err
				}
				curSolution.Parts = append(curSolution.Parts, mathbattle.Image{
					Extension: extensions[i],
					Content:   content,
				})
			}
		}
		result = append(result, curSolution)
	}

	return result, nil
}

func (r *SolutionRepository) FindOrCreate(roundID string, participantID string, problemID string) (mathbattle.Solution, error) {
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

func (r *SolutionRepository) AppendPart(ID string, item mathbattle.Image) error {
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

func (r *SolutionRepository) Update(solution mathbattle.Solution) error {
	extensions := ""
	for i := 0; i < len(solution.Parts); i++ {
		extensions = extensions + solution.Parts[i].Extension + ","
	}
	extensions = extensions[:len(extensions)-1]

	_, err := r.db.Exec("UPDATE solutions SET round_id = ?, participant_id = ?, problem_id = ?, parts = ? WHERE id = ?",
		solution.RoundID, solution.ParticipantID, solution.ProblemID, extensions, solution.ID)

	if err == sql.ErrNoRows {
		return mathbattle.ErrNotFound
	}

	return nil
}

func (r *SolutionRepository) Delete(ID string) error {
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