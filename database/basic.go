package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	mathbattle "mathbattle/models"

	_ "github.com/mattn/go-sqlite3"
)

//База задач и база участников разделены, потому что sqlite3 не способна поддерживать одновременную запись
//в одну базу данных но в разные таблицы. Может случиться так, что автор задач будет добавлять задачи в базу
//и в этом время кто-то будет регистрироваться.

type sqliteRepository struct {
	db *sql.DB
}

type SQLProblemRepository struct {
	sqliteRepository
	problemFolder string
}

type SQLParticipantRepository struct {
	sqliteRepository
}

type SQLRoundRepository struct {
	sqliteRepository
}

type SQLSolutionRepository struct {
	sqliteRepository
	solutionFolder string
}

func newSqliteRepository(dbPath string) (sqliteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return sqliteRepository{}, err
	}

	return sqliteRepository{
		db: db,
	}, nil
}

func (r *sqliteRepository) CreateFirstTime() error {
	tableCreators := []string{
		`CREATE TABLE IF NOT EXISTS participants (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			tg_chat_id VARCHAR(100),
			name VARCHAR(100),
			school VARCHAR(256),
			grade INTEGER,
			register_time DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS problems (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			sha256sum VARCHAR(64) UNIQUE,
			grade_min INTEGER,
			grade_max INTEGER,
			extension varchar(20)
		)`,
		`CREATE TABLE IF NOT EXISTS rounds (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			date_start DATETIME,
			date_end DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS round_distributions (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			round_id INTEGER,
			participant_id INTEGER,
			problem_ids TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS solutions (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			round_id INTEGER,
			participant_id INTEGER,
			problem_id INTEGER,
			parts TEXT
		)`,
	}

	for _, createStmt := range tableCreators {
		_, err := r.db.Exec(createStmt)
		if err != nil {
			return err
		}
	}

	return nil
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

func NewSQLParticipantRepository(dbPath string) (SQLParticipantRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return SQLParticipantRepository{}, err
	}

	return SQLParticipantRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func (r *SQLParticipantRepository) Store(participant mathbattle.Participant) error {
	_, err := r.db.Exec("INSERT INTO participants (tg_chat_id, name, school, grade, register_time) VALUES (?, ?, ?, ?, ?)",
		participant.TelegramID, participant.Name, participant.School, participant.Grade, participant.RegistrationTime)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLParticipantRepository) GetByID(ID string) (mathbattle.Participant, bool, error) {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return mathbattle.Participant{}, false, err
	}

	row := r.db.QueryRow("SELECT tg_chat_id, name, school, grade, register_time FROM participants WHERE id = ?", intID)
	result := mathbattle.Participant{}
	err = row.Scan(&result.TelegramID, &result.Name, &result.School, &result.Grade, &result.RegistrationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Participant{}, false, nil
		}
		return mathbattle.Participant{}, false, err
	}
	result.ID = ID

	return result, true, nil
}

func (r *SQLParticipantRepository) GetByTelegramID(telegramID string) (mathbattle.Participant, bool, error) {
	row := r.db.QueryRow("SELECT id, name, school, grade, register_time FROM participants WHERE tg_chat_id = ?", telegramID)
	var id int
	result := mathbattle.Participant{}
	err := row.Scan(&id, &result.Name, &result.School, &result.Grade, &result.RegistrationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Participant{}, false, nil
		}
		return mathbattle.Participant{}, false, err
	}
	result.ID = strconv.Itoa(id)

	return result, true, nil
}

func (r *SQLParticipantRepository) GetAll() ([]mathbattle.Participant, error) {
	rows, err := r.db.Query("SELECT id, tg_chat_id, name, school, grade, register_time FROM participants")
	if err != nil {
		return []mathbattle.Participant{}, err
	}
	defer rows.Close()

	result := []mathbattle.Participant{}
	for rows.Next() {
		var id int
		curParticipant := mathbattle.Participant{}
		err = rows.Scan(&id, &curParticipant.TelegramID, &curParticipant.Name, &curParticipant.School,
			&curParticipant.Grade, &curParticipant.Name)
		if err != nil {
			return []mathbattle.Participant{}, err
		}
		curParticipant.ID = strconv.Itoa(id)
		result = append(result, curParticipant)
	}
	return result, nil
}

func (r *SQLParticipantRepository) Delete(ID string) error {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM participants WHERE id = ?", intID)
	return err
}

func NewSQLRoundRepository(dbPath string) (SQLRoundRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return SQLRoundRepository{}, err
	}

	return SQLRoundRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func (r *SQLRoundRepository) GetDistributionForRound(roundID string) (mathbattle.RoundDistribution, error) {
	intRoundID, err := strconv.Atoi(roundID)
	if err != nil {
		return mathbattle.RoundDistribution{}, err
	}

	rows, err := r.db.Query("SELECT participant_id, problem_ids FROM round_distributions WHERE round_id=?", intRoundID)
	if err != nil {
		return mathbattle.RoundDistribution{}, err
	}
	defer rows.Close()

	var result mathbattle.RoundDistribution = make(map[string][]string)
	for rows.Next() {
		var participantID int
		var problemIDs string
		err = rows.Scan(&participantID, &problemIDs)
		if err != nil {
			return mathbattle.RoundDistribution{}, err
		}

		result[strconv.Itoa(participantID)] = strings.Split(problemIDs, ",")
	}

	return result, nil
}

func (r *SQLRoundRepository) GetAll() ([]mathbattle.Round, error) {
	rows, err := r.db.Query("SELECT id, date_start, date_end FROM rounds")
	if err != nil {
		return []mathbattle.Round{}, err
	}
	defer rows.Close()

	result := []mathbattle.Round{}
	for rows.Next() {
		var roundID int
		curRound := mathbattle.Round{}
		err = rows.Scan(&roundID, &curRound.StartDate, &curRound.EndDate)
		if err != nil {
			return []mathbattle.Round{}, err
		}
		curRound.ID = strconv.Itoa(roundID)
		distributions, err := r.GetDistributionForRound(curRound.ID)
		if err != nil {
			return []mathbattle.Round{}, err
		}
		curRound.ProblemDistribution = distributions

		result = append(result, curRound)
	}
	return result, nil
}

func (r *SQLRoundRepository) Store(round mathbattle.Round) error {
	res, err := r.db.Exec("INSERT INTO rounds (date_start, date_end) VALUES (?,?)", round.StartDate, round.EndDate)
	if err != nil {
		return err
	}

	roundID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for participantID, participantProblems := range round.ProblemDistribution {
		serializedProblems := ""
		for _, problemID := range participantProblems {
			serializedProblems = serializedProblems + problemID + ","
		}
		serializedProblems = serializedProblems[:len(serializedProblems)-1]

		intParticipantID, err := strconv.ParseInt(participantID, 10, 32)
		if err != nil {
			return err
		}

		_, err = r.db.Exec("INSERT INTO round_distributions (round_id, participant_id, problem_ids) VALUES (?,?,?)",
			roundID, intParticipantID, serializedProblems)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SQLRoundRepository) GetRunning() (mathbattle.Round, error) {
	res := r.db.QueryRow("SELECT id, date_start, date_end FROM rounds WHERE date_end = ? OR date_end <= ?",
		time.Time{}, time.Now())

	result := mathbattle.Round{}
	err := res.Scan(&result.ID, &result.StartDate, &result.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrRoundNotFound
		}
		return result, err
	}

	result.ProblemDistribution, err = r.GetDistributionForRound(result.ID)
	return result, err
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
			return result, mathbattle.ErrSolutionNotFound
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
			return mathbattle.Solution{}, mathbattle.ErrSolutionNotFound
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

	if err != mathbattle.ErrSolutionNotFound {
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
