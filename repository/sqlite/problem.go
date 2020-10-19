package sqlite

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	mathbattle "mathbattle/models"
)

type ProblemRepository struct {
	sqliteRepository
	problemFolder string
}

func NewProblemRepository(dbPath, problemPath string) (ProblemRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return ProblemRepository{}, err
	}

	return ProblemRepository{
		sqliteRepository: sqliteRepository,
		problemFolder:    problemPath,
	}, nil
}

func NewProblemRepositoryTemp(dbName, problemFolderName string) (ProblemRepository, error) {
	sqliteRepository, err := newTempSqliteRepository(dbName)
	if err != nil {
		return ProblemRepository{}, err
	}

	problemPath := filepath.Join(os.TempDir(), problemFolderName)
	os.RemoveAll(problemPath)
	err = os.Mkdir(problemPath, 0777)
	if err != nil {
		return ProblemRepository{}, err
	}

	return ProblemRepository{
		sqliteRepository: sqliteRepository,
		problemFolder:    problemPath,
	}, nil
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

	insertRes, err := r.db.Exec("INSERT INTO problems (sha256sum, grade_min, grade_max, extension) VALUES (?, ?, ?, ?)",
		problem.Sha256sum, problem.MinGrade, problem.MaxGrade, problem.Extension)
	if err != nil {
		return problem, err
	}

	id, err := insertRes.LastInsertId()
	if err != nil {
		return problem, err
	}
	problem.ID = strconv.FormatInt(id, 10)

	return problem, nil
}

func (r *ProblemRepository) GetByID(ID string) (mathbattle.Problem, error) {
	row := r.db.QueryRow("SELECT id, sha256sum, grade_min, grade_max, extension FROM problems WHERE id = ?", ID)
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
	rows, err := r.db.Query("SELECT id, sha256sum, grade_min, grade_max, extension FROM problems")
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
