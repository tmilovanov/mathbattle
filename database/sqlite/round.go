package sqlite

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	mathbattle "mathbattle/models"
)

type SQLRoundRepository struct {
	sqliteRepository
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

func NewSQLRoundRepositoryTemp(dbName string) (SQLRoundRepository, error) {
	sqliteRepository, err := newTempSqliteRepository(dbName)
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

func (r *SQLRoundRepository) Store(round mathbattle.Round) (mathbattle.Round, error) {
	result := round

	res, err := r.db.Exec("INSERT INTO rounds (date_start, date_end) VALUES (?,?)", round.StartDate, round.EndDate)
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

		intParticipantID, err := strconv.ParseInt(participantID, 10, 32)
		if err != nil {
			return result, err
		}

		_, err = r.db.Exec("INSERT INTO round_distributions (round_id, participant_id, problem_ids) VALUES (?,?,?)",
			roundID, intParticipantID, serializedProblems)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (r *SQLRoundRepository) GetRunning() (mathbattle.Round, error) {
	res := r.db.QueryRow("SELECT id, date_start, date_end FROM rounds WHERE date_end = ? OR date_end <= ?",
		time.Time{}, time.Now())

	result := mathbattle.Round{}
	err := res.Scan(&result.ID, &result.StartDate, &result.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

	result.ProblemDistribution, err = r.GetDistributionForRound(result.ID)
	return result, err
}
