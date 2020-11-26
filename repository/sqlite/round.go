package sqlite

import (
	"database/sql"
	"encoding/json"
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

func serializeProblemIDs(problemIDs []string) string {
	return strings.Join(problemIDs, ",")
}

func deserializeProblemIDs(serializedProblemIDs string) []string {
	return strings.Split(serializedProblemIDs, ",")
}

func (r *RoundRepository) ProblemDistributionStore(roundID string, rd mathbattle.RoundDistribution) error {
	for participantID, participantProblems := range rd {
		_, err := r.db.Exec("INSERT INTO rounds_problems_distributions (round_id, participant_id, problems_ids) VALUES (?,?,?)",
			roundID, participantID, serializeProblemIDs(participantProblems))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RoundRepository) ProblemDistributionGet(roundID string) (mathbattle.RoundDistribution, error) {
	rows, err := r.db.Query("SELECT participant_id, problems_ids FROM rounds_problems_distributions WHERE round_id=?", roundID)
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

		result[strconv.Itoa(participantID)] = deserializeProblemIDs(problemIDs)
	}

	return result, nil
}

func (r *RoundRepository) ProblemDistributionUpdate(roundID string, rd mathbattle.RoundDistribution) error {
	for participantID, participantProblems := range rd {
		_, err := r.db.Exec("REPLACE INTO rounds_problems_distributions (round_id, participant_id, problems_ids) VALUES (?,?,?)",
			roundID, participantID, serializeProblemIDs(participantProblems))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RoundRepository) ProblemDistributionDelete(roundID string) error {
	_, err := r.db.Exec("DELETE FROM rounds_problems_distributions WHERE round_id = ?", roundID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mathbattle.ErrNotFound
		}
		return err
	}
	return nil
}

type ReviewDistribution struct {
	BetweenParticipants map[string][]string `json:"between_participants"`
	ToOrganizers        []string            `json:"to_organizers"`
}

func serializeReviewDistribution(rd mathbattle.ReviewDistribution) (string, error) {
	localRd := ReviewDistribution{
		BetweenParticipants: rd.BetweenParticipants,
		ToOrganizers:        rd.ToOrganizers,
	}

	serialized, err := json.Marshal(localRd)
	return string(serialized), err
}

func deserializeReviewDistribution(input string) (mathbattle.ReviewDistribution, error) {
	var rd ReviewDistribution
	err := json.Unmarshal([]byte(input), &rd)
	if err != nil {
		return mathbattle.ReviewDistribution{}, err
	}

	return mathbattle.ReviewDistribution{
		BetweenParticipants: rd.BetweenParticipants,
		ToOrganizers:        rd.ToOrganizers,
	}, nil
}

func (r *RoundRepository) SolutionDistributionStore(roundID string, rd mathbattle.ReviewDistribution) error {
	serialized, err := serializeReviewDistribution(rd)
	if err != nil {
		return nil
	}

	_, err = r.db.Exec("INSERT INTO rounds_solutions_distributions (round_id, distribution) VALUES (?,?)",
		roundID, serialized)
	return err
}

func (r *RoundRepository) SolutionDistributionGet(roundID string) (mathbattle.ReviewDistribution, error) {
	res := r.db.QueryRow("SELECT distribution FROM rounds_solutions_distributions WHERE round_id=?", roundID)
	var serialized string
	err := res.Scan(&serialized)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.ReviewDistribution{}, mathbattle.ErrNotFound
		}
		return mathbattle.ReviewDistribution{}, err
	}

	return deserializeReviewDistribution(serialized)
}

func (r *RoundRepository) SolutionDistributionUpdate(roundID string, rd mathbattle.ReviewDistribution) error {
	serialized, err := serializeReviewDistribution(rd)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("UPDATE rounds_solutions_distributions SET distribution = ? WHERE round_id = ?", serialized, roundID)
	return err
}

func (r *RoundRepository) SolutionDistributionDelete(roundID string) error {
	_, err := r.db.Exec("DELETE FROM rounds_solutions_distributions WHERE round_id = ?", roundID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mathbattle.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *RoundRepository) Store(round mathbattle.Round) (mathbattle.Round, error) {
	res, err := r.db.Exec("INSERT INTO rounds (solve_start, solve_end, review_start, review_end) VALUES (?,?,?,?)",
		round.GetSolveStartDate(), round.GetSolveEndDate(), round.GetReviewStartDate(), round.GetReviewEndDate())
	if err != nil {
		return round, err
	}

	roundID, err := res.LastInsertId()
	if err != nil {
		return round, err
	}
	round.ID = strconv.FormatInt(roundID, 10)

	err = r.ProblemDistributionStore(round.ID, round.ProblemDistribution)
	if err != nil {
		return round, err
	}

	err = r.SolutionDistributionStore(round.ID, round.ReviewDistribution)

	return round, err
}

func (r *RoundRepository) Get(ID string) (mathbattle.Round, error) {
	result := mathbattle.Round{ID: ID}

	res := r.db.QueryRow("SELECT solve_start, solve_end, review_start, review_end FROM rounds WHERE id = ?", ID)
	var solveStartDate time.Time
	var solveEndDate time.Time
	var reviewStartDate time.Time
	var reviewEndDate time.Time
	err := res.Scan(&solveStartDate, &solveEndDate, &reviewStartDate, &reviewEndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

	result.SetSolveStartDate(solveStartDate)
	result.SetSolveEndDate(solveEndDate)
	result.SetReviewStartDate(reviewStartDate)
	result.SetReviewEndDate(reviewEndDate)

	result.ProblemDistribution, err = r.ProblemDistributionGet(result.ID)
	if err != nil {
		return result, err
	}

	result.ReviewDistribution, err = r.SolutionDistributionGet(result.ID)
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
		var solveStartDate time.Time
		var solveEndDate time.Time
		var reviewStartDate time.Time
		var reviewEndDate time.Time
		err = rows.Scan(&roundID, &solveStartDate, &solveEndDate, &reviewStartDate, &reviewEndDate)
		if err != nil {
			return []mathbattle.Round{}, err
		}
		curRound.ID = strconv.Itoa(roundID)
		curRound.ProblemDistribution, err = r.ProblemDistributionGet(curRound.ID)
		curRound.SetSolveStartDate(solveStartDate)
		curRound.SetSolveEndDate(solveEndDate)
		curRound.SetReviewStartDate(reviewStartDate)
		curRound.SetReviewEndDate(reviewEndDate)
		if err != nil {
			return []mathbattle.Round{}, err
		}
		curRound.ReviewDistribution, err = r.SolutionDistributionGet(curRound.ID)
		if err != nil {
			return []mathbattle.Round{}, err
		}

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
		round.GetSolveStartDate(), round.GetSolveEndDate(), round.GetReviewStartDate(), round.GetReviewEndDate(), round.ID)
	if err != nil {
		return err
	}

	err = r.ProblemDistributionUpdate(round.ID, round.ProblemDistribution)
	if err != nil {
		return err
	}

	return r.SolutionDistributionUpdate(round.ID, round.ReviewDistribution)
}

func (r *RoundRepository) Delete(ID string) error {
	if err := r.ProblemDistributionDelete(ID); err != nil {
		return err
	}

	if err := r.SolutionDistributionDelete(ID); err != nil {
		return err
	}

	_, err := r.db.Exec("DELETE FROM rounds WHERE id = ?", ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mathbattle.ErrNotFound
		}
		return err
	}

	return nil
}
