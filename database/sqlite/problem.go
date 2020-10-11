package sqlite

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	mathbattle "mathbattle/models"
)

type SQLProblemRepository struct {
	sqliteRepository
	problemFolder string
}

func NewSQLProblemRepository(dbPath, problemPath string) (SQLProblemRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return SQLProblemRepository{}, err
	}

	return SQLProblemRepository{
		sqliteRepository: sqliteRepository,
		problemFolder:    problemPath,
	}, nil
}

func (r *SQLProblemRepository) getFilePathFromProblem(problem mathbattle.Problem) string {
	return filepath.Join(r.problemFolder, fmt.Sprintf("%d_%d_%s%s",
		problem.MinGrade, problem.MaxGrade, problem.ID, problem.Extension))
}

func (r *SQLProblemRepository) Store(problem mathbattle.Problem) error {
	err := ioutil.WriteFile(r.getFilePathFromProblem(problem), problem.Content, 0666)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("REPLACE INTO problems (sha256sum, grade_min, grade_max, extension) VALUES (?, ?, ?, ?)",
		problem.ID, problem.MinGrade, problem.MaxGrade, problem.Extension)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLProblemRepository) GetByID(ID string) (mathbattle.Problem, error) {
	row := r.db.QueryRow("SELECT sha256sum, grade_min, grade_max, extension FROM problems WHERE sha256sum = ?", ID)
	result := mathbattle.Problem{}
	err := row.Scan(&result.ID, &result.MinGrade, &result.MaxGrade, &result.Extension)
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

func (r *SQLProblemRepository) GetAll() ([]mathbattle.Problem, error) {
	rows, err := r.db.Query("SELECT sha256sum, grade_min, grade_max, extension FROM problems")
	if err != nil {
		return []mathbattle.Problem{}, err
	}
	defer rows.Close()

	result := []mathbattle.Problem{}
	for rows.Next() {
		curProblem := mathbattle.Problem{}
		err = rows.Scan(&curProblem.ID, &curProblem.MinGrade, &curProblem.MaxGrade, &curProblem.Extension)
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
