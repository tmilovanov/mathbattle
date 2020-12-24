package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"mathbattle/models/mathbattle"
)

type RoundRepository struct {
	sqliteRepository
}

func NewRoundRepository(dbPath string) (*RoundRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return nil, err
	}

	result := &RoundRepository{
		sqliteRepository: sqliteRepository,
	}

	if err := result.CreateTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RoundRepository) CreateTable() error {
	createStmt := `CREATE TABLE IF NOT EXISTS rounds (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			solve_start DATETIME,
			solve_end DATETIME,
			review_start DATETIME,
			review_end DATETIME,
			problems_distribution TEXT,
			solutions_distribution TEXT
		)`
	_, err := r.db.Exec(createStmt)
	return err
}

type ProblemDescriptor struct {
	Caption   string `json:"caption"`
	ProblemID string `json:"id"`
}

type RoundDistribution struct {
	Distribution map[string][]ProblemDescriptor `json:"problems_distribution"`
}

func serializeProblemsDistribution(rd mathbattle.RoundDistribution) (string, error) {
	localRoundDistribution := RoundDistribution{
		Distribution: make(map[string][]ProblemDescriptor),
	}
	for key, value := range rd {
		for _, desc := range value {
			localRoundDistribution.Distribution[key] = append(localRoundDistribution.Distribution[key], ProblemDescriptor{
				Caption:   desc.Caption,
				ProblemID: desc.ProblemID,
			})
		}
	}

	serliazed, err := json.Marshal(&localRoundDistribution)
	return string(serliazed), err
}

func deserializeProblemsDistribution(input string) (mathbattle.RoundDistribution, error) {
	var rd RoundDistribution
	err := json.Unmarshal([]byte(input), &rd)
	if err != nil {
		return mathbattle.RoundDistribution{}, err
	}

	result := mathbattle.RoundDistribution{}
	for key, value := range rd.Distribution {
		for _, desc := range value {
			result[key] = append(result[key], mathbattle.ProblemDescriptor{
				Caption:   desc.Caption,
				ProblemID: desc.ProblemID,
			})
		}
	}
	return result, nil
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

func (r *RoundRepository) Store(round mathbattle.Round) (mathbattle.Round, error) {
	serializedRoundDistribution, err := serializeProblemsDistribution(round.ProblemDistribution)
	if err != nil {
		return round, err
	}
	serializedSolutionDistribution, err := serializeReviewDistribution(round.ReviewDistribution)
	if err != nil {
		return round, err
	}
	res, err := r.db.Exec(`INSERT INTO rounds (solve_start, solve_end, review_start, review_end,
		problems_distribution, solutions_distribution) VALUES (?,?,?,?,?,?)`, round.GetSolveStartDate(), round.GetSolveEndDate(),
		round.GetReviewStartDate(), round.GetReviewEndDate(), serializedRoundDistribution, serializedSolutionDistribution)
	if err != nil {
		return round, err
	}

	roundID, err := res.LastInsertId()
	if err != nil {
		return round, err
	}
	round.ID = strconv.FormatInt(roundID, 10)

	return round, err
}

func (r *RoundRepository) Get(ID string) (mathbattle.Round, error) {
	result := mathbattle.Round{ID: ID}

	res := r.db.QueryRow(`SELECT solve_start, solve_end, review_start, review_end, 
	problems_distribution, solutions_distribution FROM rounds WHERE id = ?`, ID)
	var solveStartDate time.Time
	var solveEndDate time.Time
	var reviewStartDate time.Time
	var reviewEndDate time.Time
	var problemsDistribution string
	var solutionsDistribution string
	err := res.Scan(&solveStartDate, &solveEndDate, &reviewStartDate, &reviewEndDate,
		&problemsDistribution, &solutionsDistribution)
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
	result.ProblemDistribution, err = deserializeProblemsDistribution(problemsDistribution)
	if err != nil {
		return result, err
	}
	result.ReviewDistribution, err = deserializeReviewDistribution(solutionsDistribution)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *RoundRepository) GetAll() ([]mathbattle.Round, error) {
	rows, err := r.db.Query(`SELECT id, solve_start, solve_end, review_start, review_end,
		problems_distribution, solutions_distribution FROM rounds`)
	if err != nil {
		return []mathbattle.Round{}, err
	}
	defer rows.Close()

	result := []mathbattle.Round{}
	for rows.Next() {
		curRound := mathbattle.Round{}
		var roundID int
		var solveStartDate time.Time
		var solveEndDate time.Time
		var reviewStartDate time.Time
		var reviewEndDate time.Time
		var problemsDistribution string
		var solutionsDistribution string
		err = rows.Scan(&roundID, &solveStartDate, &solveEndDate, &reviewStartDate, &reviewEndDate,
			&problemsDistribution, &solutionsDistribution)
		if err != nil {
			return []mathbattle.Round{}, err
		}
		curRound.ID = strconv.Itoa(roundID)
		curRound.SetSolveStartDate(solveStartDate)
		curRound.SetSolveEndDate(solveEndDate)
		curRound.SetReviewStartDate(reviewStartDate)
		curRound.SetReviewEndDate(reviewEndDate)
		curRound.ProblemDistribution, err = deserializeProblemsDistribution(problemsDistribution)
		if err != nil {
			return result, err
		}
		curRound.ReviewDistribution, err = deserializeReviewDistribution(solutionsDistribution)
		if err != nil {
			return result, err
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
		time.Time{}, time.Now().Round(0).UTC())

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
		time.Now().Round(0).UTC(), time.Time{})

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
	res := r.db.QueryRow(`
		SELECT id FROM rounds WHERE
			solve_end <= ? AND
			(review_start != ? AND review_start <= ?) AND
			(review_end = ? OR review_end >= ?)`,
		time.Now().Round(0).UTC(),
		time.Time{}, time.Now().Round(0).UTC(),
		time.Time{}, time.Now().Round(0).UTC())

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

func (r *RoundRepository) GetLast() (mathbattle.Round, error) {
	res := r.db.QueryRow("SELECT id FROM rounds ORDER BY ID DESC LIMIT 1")

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
	serializedRoundDistribution, err := serializeProblemsDistribution(round.ProblemDistribution)
	if err != nil {
		return err
	}
	serializedSolutionDistribution, err := serializeReviewDistribution(round.ReviewDistribution)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE rounds SET solve_start = ?, solve_end = ?, review_start = ?, review_end = ?,
	problems_distribution = ?, solutions_distribution = ? WHERE id = ?`, round.GetSolveStartDate(), round.GetSolveEndDate(),
		round.GetReviewStartDate(), round.GetReviewEndDate(), serializedRoundDistribution, serializedSolutionDistribution, round.ID)
	return err
}

func (r *RoundRepository) Delete(ID string) error {
	_, err := r.db.Exec("DELETE FROM rounds WHERE id = ?", ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mathbattle.ErrNotFound
		}
		return err
	}

	return nil
}
