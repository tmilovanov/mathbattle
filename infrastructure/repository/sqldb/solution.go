package sqldb

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
	sqlRepository
	solutionFolder string
}

func NewSolutionRepository(dbType, connectionString, solutionPath string) (*SolutionRepository, error) {
	sqlRepository, err := newSqlRepository(dbType, connectionString)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		if err := os.MkdirAll(solutionPath, 0777); err != nil {
			return nil, err
		}
	}

	result := &SolutionRepository{
		sqlRepository:  sqlRepository,
		solutionFolder: solutionPath,
	}

	if err := result.CreateTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *SolutionRepository) CreateTable() error {
	var createStmt string

	switch r.dbType {
	case "sqlite3":
		createStmt = `CREATE TABLE IF NOT EXISTS solutions (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			round_id INTEGER,
			participant_id INTEGER,
			problem_id INTEGER,
			juri_comment TEXT,
			mark INTEGER,
			parts TEXT
		)`
	case "postgres":
		createStmt = `CREATE TABLE IF NOT EXISTS solutions (
			id SERIAL UNIQUE,
			round_id INTEGER,
			participant_id INTEGER,
			problem_id INTEGER,
			juri_comment TEXT,
			mark INTEGER,
			parts TEXT
		)`
	}

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

	switch r.dbType {
	case "sqlite3":
		res, err := r.db.Exec("INSERT INTO solutions (round_id, participant_id, problem_id, juri_comment, mark, parts) VALUES (?,?,?,?,?,?)",
			solution.RoundID, solution.ParticipantID, solution.ProblemID, solution.JuriComment, extensions)
		if err != nil {
			return result, err
		}

		newID, err := res.LastInsertId()
		if err != nil {
			return result, err
		}
		result.ID = strconv.FormatInt(newID, 10)

		return result, nil
	case "postgres":
		query := "INSERT INTO solutions (round_id, participant_id, problem_id, juri_comment, mark, parts) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
		stmt, err := r.db.Prepare(query)
		if err != nil {
			return result, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(solution.RoundID, solution.ParticipantID, solution.ProblemID, solution.JuriComment, solution.Mark, extensions).Scan(&result.ID)
		if err != nil {
			return result, err
		}

		return result, nil
	default:
		return result, fmt.Errorf("Unknown dbtype")
	}
}

func (r *SolutionRepository) getManyWhere(whereStr string, whereArgs ...interface{}) ([]mathbattle.Solution, error) {
	result := []mathbattle.Solution{}

	rows, err := r.db.Query(`
	SELECT id, round_id, participant_id, problem_id, juri_comment, mark, parts
	FROM solutions
	WHERE `+whereStr, whereArgs...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var cur mathbattle.Solution
		var partsExtensions string
		err := rows.Scan(&cur.ID, &cur.RoundID, &cur.ParticipantID, &cur.ProblemID, &cur.JuriComment, &cur.Mark, &partsExtensions)
		if err != nil {
			return result, err
		}

		if len(partsExtensions) != 0 {
			extensions := strings.Split(partsExtensions, ",")
			for i := 0; i < len(extensions); i++ {
				partPath := r.getPartPath(cur, i, extensions[i])
				content, err := ioutil.ReadFile(partPath)
				if err != nil {
					return result, err
				}
				cur.Parts = append(cur.Parts, mathbattle.Image{
					Extension: extensions[i],
					Content:   content,
				})
			}
		}

		result = append(result, cur)
	}

	return result, nil
}

func (r *SolutionRepository) getOneWhere(whereStr string, whereArgs ...interface{}) (mathbattle.Solution, error) {
	res, err := r.getManyWhere(whereStr, whereArgs...)
	if err != nil {
		return mathbattle.Solution{}, err
	}

	if len(res) == 0 {
		return mathbattle.Solution{}, mathbattle.ErrNotFound
	}

	return res[0], nil
}

func (r *SolutionRepository) Get(ID string) (mathbattle.Solution, error) {
	return r.getOneWhere("id = $1", ID)
}

func (r *SolutionRepository) Find(roundID string, participantID string, problemID string) (mathbattle.Solution, error) {
	return r.getOneWhere("round_id = $1 AND participant_id = $2 AND problem_id = $3",
		roundID, participantID, problemID)
}

func (r *SolutionRepository) FindMany(roundID string, participantID string, problemID string) ([]mathbattle.Solution, error) {
	whereClause, whereArgs := joinWhereOmitEmpty([]whereDescriptor{
		{"round_id", roundID},
		{"participant_id", participantID},
		{"problem_id", problemID},
	})
	return r.getManyWhere(whereClause, whereArgs...)
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
	_, err = r.db.Exec("UPDATE solutions SET parts = $1 WHERE id = $2", extensions, intID)

	return err
}

func (r *SolutionRepository) Update(solution mathbattle.Solution) error {
	extensions := ""
	for i := 0; i < len(solution.Parts); i++ {
		extensions = extensions + solution.Parts[i].Extension + ","
	}
	extensions = extensions[:len(extensions)-1]

	_, err := r.db.Exec(`
	UPDATE solutions
	SET round_id = $1, participant_id = $2, problem_id = $3, juri_comment=$4, mark=$5, parts = $6
	WHERE id = $7`,
		solution.RoundID, solution.ParticipantID, solution.ProblemID, solution.JuriComment, solution.Mark, extensions, solution.ID)

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

	_, err = r.db.Exec("DELETE FROM solutions WHERE id = $1", solution.ID)

	return err
}
