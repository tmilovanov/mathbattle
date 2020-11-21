package sqlite

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	mathbattle "mathbattle/models"
)

type RoundRepository struct {
	sqliteRepository
}

func NewRoundRepository(dbPath string) (RoundRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return RoundRepository{}, err
	}

	return RoundRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func NewRoundRepositoryTemp(dbName string) (RoundRepository, error) {
	sqliteRepository, err := newTempSqliteRepository(dbName)
	if err != nil {
		return RoundRepository{}, err
	}

	return RoundRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func (r *RoundRepository) ProblemDistributionGet(roundID string) (mathbattle.RoundDistribution, error) {
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

func (r *RoundRepository) ProblemDistributionUpdate(roundID string, rd mathbattle.RoundDistribution) error {
	return nil
}

func (r *RoundRepository) Store(round mathbattle.Round) (mathbattle.Round, error) {
	result := round

	res, err := r.db.Exec("INSERT INTO rounds (solve_start, solve_end, review_start, review_end) VALUES (?,?,?,?)",
		round.SolveStartDate, round.SolveEndDate, round.ReviewStartDate, round.ReviewEndDate)
	if err != nil {
		return result, err
	}

	roundID, err := res.LastInsertId()
	if err != nil {
		return result, err
	}
	result.ID = strconv.FormatInt(roundID, 10)

	for participantID, participantProblems := range round.ProblemDistribution {
		serializedProblems := ""
		for _, problemID := range participantProblems {
			serializedProblems = serializedProblems + problemID + ","
		}
		serializedProblems = serializedProblems[:len(serializedProblems)-1]

		_, err = r.db.Exec("INSERT INTO round_distributions (round_id, participant_id, problem_ids) VALUES (?,?,?)",
			roundID, participantID, serializedProblems)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (r *RoundRepository) Get(ID string) (mathbattle.Round, error) {
	result := mathbattle.Round{ID: ID}

	intID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return result, err
	}

	res := r.db.QueryRow("SELECT solve_start, solve_end, review_start, review_end FROM rounds WHERE id = ?", intID)
	err = res.Scan(&result.SolveStartDate, &result.SolveEndDate, &result.ReviewStartDate, &result.ReviewEndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, sql.ErrConnDone
		}
		return result, err
	}

	result.ProblemDistribution, err = r.ProblemDistributionGet(result.ID)
	return result, err
}

func (r *RoundRepository) GetAll() ([]mathbattle.Round, error) {
	rows, err := r.db.Query("SELECT id, solve_start, solve_end, review_start, review_end FROM rounds")
	if err != nil {
		return []mathbattle.Round{}, err
	}
	defer rows.Close()

	result := []mathbattle.Round{}
	for rows.Next() {
		var roundID int
		curRound := mathbattle.Round{}
		err = rows.Scan(&roundID, &curRound.SolveStartDate, &curRound.SolveEndDate, &curRound.ReviewStartDate, &curRound.ReviewEndDate)
		if err != nil {
			return []mathbattle.Round{}, err
		}
		curRound.ID = strconv.Itoa(roundID)
		distributions, err := r.ProblemDistributionGet(curRound.ID)
		if err != nil {
			return []mathbattle.Round{}, err
		}
		curRound.ProblemDistribution = distributions

		result = append(result, curRound)
	}
	return result, nil
}

func (r *RoundRepository) GetRunning() (mathbattle.Round, error) {
	round, err := r.GetSolveRunning()
	if err == nil {
		return round, nil
	}

	if err != mathbattle.ErrNotFound {
		return mathbattle.Round{}, err
	}

	round, err = r.GetReviewPending()
	if err == nil {
		return round, nil
	}

	if err != mathbattle.ErrNotFound {
		return mathbattle.Round{}, err
	}

	return r.GetReviewRunning()
}

func (r *RoundRepository) GetSolveRunning() (mathbattle.Round, error) {
	res := r.db.QueryRow("SELECT id FROM rounds WHERE solve_end = ? OR solve_end >= ?",
		time.Time{}, time.Now())

	var ID string
	err := res.Scan(&ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Round{}, mathbattle.ErrNotFound
		}
		return mathbattle.Round{}, err
	}

	return r.Get(ID)
}

func (r *RoundRepository) GetReviewPending() (mathbattle.Round, error) {
	res := r.db.QueryRow("SELECT id FROM rounds WHERE solve_end <= ? AND review_start = ?",
		time.Now(), time.Time{})

	var ID string
	err := res.Scan(&ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Round{}, mathbattle.ErrNotFound
		}
		return mathbattle.Round{}, err
	}

	return r.Get(ID)
}

func (r *RoundRepository) GetReviewRunning() (mathbattle.Round, error) {
	res := r.db.QueryRow("SELECT id FROM rounds WHERE solve_end <= ? AND (review_end = ? OR review_end >= ?)",
		time.Now(), time.Time{}, time.Now())

	var ID string
	err := res.Scan(&ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Round{}, mathbattle.ErrNotFound
		}
		return mathbattle.Round{}, err
	}

	return r.Get(ID)
}

func (r *RoundRepository) Update(round mathbattle.Round) error {
	_, err := r.db.Exec("UPDATE rounds SET solve_start = ?, solve_end = ?, review_start = ?, review_end = ? WHERE id = ?",
		round.SolveStartDate, round.SolveEndDate, round.ReviewStartDate, round.ReviewEndDate, round.ID)
	if err != nil {
		return err
	}

	return r.ProblemDistributionUpdate(round.ID, round.ProblemDistribution)
}

func (r *RoundRepository) Delete(ID string) error {
	intID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM round_distributions WHERE round_id = ?", intID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mathbattle.ErrNotFound
		}
		return err
	}

	_, err = r.db.Exec("DELETE FROM rounds WHERE id = ?", intID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mathbattle.ErrNotFound
		}
		return err
	}

	return nil
}
