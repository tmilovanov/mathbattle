package sqldb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"mathbattle/models/mathbattle"
)

type RoundRepository struct {
	sqlRepository
}

func NewRoundRepository(dbType, connectionString string) (*RoundRepository, error) {
	sqlRepository, err := newSqlRepository(dbType, connectionString)
	if err != nil {
		return nil, err
	}

	result := &RoundRepository{
		sqlRepository: sqlRepository,
	}

	if err := result.CreateTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RoundRepository) CreateTable() error {
	var createStmt string

	switch r.dbType {
	case "sqlite3":
		createStmt = `CREATE TABLE IF NOT EXISTS rounds (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			solve_start DATETIME,
			solve_end DATETIME,
			review_start DATETIME,
			review_end DATETIME,
			problems_distribution TEXT,
			solutions_distribution TEXT
		)`
	case "postgres":
		createStmt = `CREATE TABLE IF NOT EXISTS rounds (
			id SERIAL UNIQUE,
			solve_start TIMESTAMP,
			solve_end TIMESTAMP,
			review_start TIMESTAMP,
			review_end TIMESTAMP,
			problems_distribution TEXT,
			solutions_distribution TEXT
		)`
	}

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

	switch r.dbType {
	case "sqlite3":
		res, err := r.db.Exec(`INSERT INTO rounds (solve_start, solve_end, review_start, review_end,
		problems_distribution, solutions_distribution) VALUES ($1,$2,$3,$4,$5,$6)`,
			round.GetSolveStartDate(), round.GetSolveEndDate(),
			round.GetReviewStartDate(), round.GetReviewEndDate(),
			serializedRoundDistribution, serializedSolutionDistribution)
		if err != nil {
			return round, err
		}

		roundID, err := res.LastInsertId()
		if err != nil {
			return round, err
		}
		round.ID = strconv.FormatInt(roundID, 10)
	case "postgres":
		query := `INSERT INTO rounds (solve_start, solve_end, review_start, review_end,
		problems_distribution, solutions_distribution) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
		stmt, err := r.db.Prepare(query)
		if err != nil {
			return round, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(round.GetSolveStartDate(), round.GetSolveEndDate(),
			round.GetReviewStartDate(), round.GetReviewEndDate(),
			serializedRoundDistribution, serializedSolutionDistribution).Scan(&round.ID)
		if err != nil {
			return round, err
		}

	}

	return round, err
}

func (r *RoundRepository) getWhere(whereStr string, whereArgs ...interface{}) (mathbattle.Round, error) {
	result := mathbattle.Round{}
	res := r.db.QueryRow(`SELECT id, solve_start, solve_end, review_start, review_end, 
	problems_distribution, solutions_distribution FROM rounds WHERE `+whereStr, whereArgs...)
	var problemsDistribution string
	var solutionsDistribution string
	err := res.Scan(&result.ID, &result.SolveStartDate, &result.SolveEndDate,
		&result.ReviewStartDate, &result.ReviewEndDate,
		&problemsDistribution, &solutionsDistribution)
	result.SetSolveStartDate(result.SolveStartDate)
	result.SetSolveEndDate(result.SolveEndDate)
	result.SetReviewStartDate(result.ReviewStartDate)
	result.SetReviewEndDate(result.ReviewEndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

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

func (r *RoundRepository) Get(ID string) (mathbattle.Round, error) {
	return r.getWhere("id = $1", ID)
}

func (r *RoundRepository) GetAll() ([]mathbattle.Round, error) {
	result := []mathbattle.Round{}
	rows, err := r.db.Query("SELECT id FROM rounds")
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var curID string
		err = rows.Scan(&curID)
		if err != nil {
			return result, err
		}

		cur, err := r.Get(curID)
		if err != nil {
			return result, err
		}

		result = append(result, cur)
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
	return r.getWhere("solve_end = $1 OR solve_end >= $2",
		time.Time{}, time.Now().Round(0).UTC())
}

func (r *RoundRepository) GetReviewPending() (mathbattle.Round, error) {
	return r.getWhere("solve_end <= $1 AND review_start = $2",
		time.Now().Round(0).UTC(), time.Time{})
}

func (r *RoundRepository) GetReviewRunning() (mathbattle.Round, error) {
	return r.getWhere(`solve_end <= $1 AND
			(review_start != $2 AND review_start <= $3) AND
			(review_end = $4 OR review_end >= $5)`,
		time.Now().Round(0).UTC(),
		time.Time{}, time.Now().Round(0).UTC(),
		time.Time{}, time.Now().Round(0).UTC())
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
	_, err = r.db.Exec(`UPDATE rounds SET solve_start = $1, solve_end = $2, review_start = $3, review_end = $4,
	problems_distribution = $5, solutions_distribution = $6 WHERE id = $7`, round.GetSolveStartDate(), round.GetSolveEndDate(),
		round.GetReviewStartDate(), round.GetReviewEndDate(), serializedRoundDistribution, serializedSolutionDistribution, round.ID)
	return err
}

func (r *RoundRepository) Delete(ID string) error {
	_, err := r.db.Exec("DELETE FROM rounds WHERE id = $1", ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mathbattle.ErrNotFound
		}
		return err
	}

	return nil
}
