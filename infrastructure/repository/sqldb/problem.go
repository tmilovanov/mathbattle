package sqldb

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"mathbattle/models/mathbattle"
)

type ProblemRepository struct {
	sqlRepository
	problemFolder string
}

func NewProblemRepository(dbType, connectionString, problemPath string) (*ProblemRepository, error) {
	sqlRepository, err := newSqlRepository(dbType, connectionString)
	if err != nil {
		log.Printf("Fld to get sql rep, err: %v", err)
		return nil, err
	}

	if _, err := os.Stat(problemPath); os.IsNotExist(err) {
		if err := os.MkdirAll(problemPath, 0777); err != nil {
			log.Printf("Fld to create %v rep, err: %v", problemPath, err)
			return nil, err
		}
	}

	result := &ProblemRepository{
		sqlRepository: sqlRepository,
		problemFolder: problemPath,
	}

	if err := result.CreateTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ProblemRepository) CreateTable() error {
	var createStmt string

	switch r.dbType {
	case "sqlite3":
		createStmt = `CREATE TABLE IF NOT EXISTS problems (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			sha256sum VARCHAR(64) UNIQUE,
			grade_min INTEGER,
			grade_max INTEGER,
			extension varchar(20)
		)`
	case "postgres":
		createStmt = `CREATE TABLE IF NOT EXISTS problems (
			id SERIAL UNIQUE,
			sha256sum VARCHAR(64) UNIQUE,
			grade_min INTEGER,
			grade_max INTEGER,
			extension varchar(20)
		)`
	default:
		return fmt.Errorf("Unsupported database type")
	}

	_, err := r.db.Exec(createStmt)
	return err
}

func (r *ProblemRepository) getFilePathFromProblem(problem mathbattle.Problem) string {
	return filepath.Join(r.problemFolder, fmt.Sprintf("%d_%d_%s%s",
		problem.MinGrade, problem.MaxGrade, problem.Sha256sum, problem.Extension))
}

func (r *ProblemRepository) Store(problem mathbattle.Problem) (mathbattle.Problem, error) {
	h := sha256.New()
	if _, err := io.Copy(h, bytes.NewReader(problem.Content)); err != nil {
		return problem, err
	}
	problem.Sha256sum = hex.EncodeToString(h.Sum(nil))

	err := ioutil.WriteFile(r.getFilePathFromProblem(problem), problem.Content, 0666)
	if err != nil {
		return problem, err
	}

	switch r.dbType {
	case "sqlite3":
		insertRes, err := r.db.Exec("INSERT INTO problems (sha256sum, grade_min, grade_max, extension) VALUES ($1, $2, $3, $4)",
			problem.Sha256sum, problem.MinGrade, problem.MaxGrade, problem.Extension)
		if err != nil {
			return problem, err
		}

		id, err := insertRes.LastInsertId()
		if err != nil {
			return problem, err
		}

		problem.ID = strconv.FormatInt(id, 10)
	case "postgres":
		query := "INSERT INTO problems (sha256sum, grade_min, grade_max, extension) VALUES ($1, $2, $3, $4) RETURNING id"
		stmt, err := r.db.Prepare(query)
		if err != nil {
			return problem, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(problem.Sha256sum, problem.MinGrade, problem.MaxGrade, problem.Extension).Scan(&problem.ID)
		if err != nil {
			return problem, err
		}
	}

	return problem, nil
}

func (r *ProblemRepository) GetByID(ID string) (mathbattle.Problem, error) {
	row := r.db.QueryRow("SELECT id, sha256sum, grade_min, grade_max, extension FROM problems WHERE id = $1", ID)
	result := mathbattle.Problem{}
	err := row.Scan(&result.ID, &result.Sha256sum, &result.MinGrade, &result.MaxGrade, &result.Extension)
	if err != nil {
		return mathbattle.Problem{}, err
	}
	result.ID = ID
	content, err := ioutil.ReadFile(r.getFilePathFromProblem(result))
	if err != nil {
		return mathbattle.Problem{}, err
	}
	result.Content = content
	return result, nil
}

func (r *ProblemRepository) GetAll() ([]mathbattle.Problem, error) {
	rows, err := r.db.Query("SELECT id, sha256sum, grade_min, grade_max, extension FROM problems ORDER BY id")
	if err != nil {
		return []mathbattle.Problem{}, err
	}
	defer rows.Close()

	result := []mathbattle.Problem{}
	for rows.Next() {
		curProblem := mathbattle.Problem{}
		err = rows.Scan(&curProblem.ID, &curProblem.Sha256sum, &curProblem.MinGrade, &curProblem.MaxGrade, &curProblem.Extension)
		if err != nil {
			return []mathbattle.Problem{}, err
		}

		content, err := ioutil.ReadFile(r.getFilePathFromProblem(curProblem))
		if err != nil {
			return []mathbattle.Problem{}, err
		}
		curProblem.Content = content

		result = append(result, curProblem)
	}

	return result, nil
}
