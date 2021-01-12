package sqldb

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"mathbattle/models/mathbattle"
)

type ReviewRepository struct {
	sqlRepository
}

func NewReviewRepository(dbType, connectionString string) (*ReviewRepository, error) {
	sqlRepository, err := newSqlRepository(dbType, connectionString)
	if err != nil {
		return nil, err
	}

	result := &ReviewRepository{
		sqlRepository: sqlRepository,
	}

	if err := result.CreateTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ReviewRepository) CreateTable() error {
	var createStmt string

	switch r.dbType {
	case "sqlite3":
		createStmt = `CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			reviewer_id INTEGER,
			solution_id INTEGER,
			content TEXT
		)`
	case "postgres":
		createStmt = `CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL UNIQUE,
			reviewer_id INTEGER,
			solution_id INTEGER,
			content TEXT
		)`
	}
	_, err := r.db.Exec(createStmt)
	return err
}

func (r *ReviewRepository) Store(review mathbattle.Review) (mathbattle.Review, error) {
	result := review

	switch r.dbType {
	case "sqlite3":
		res, err := r.db.Exec("INSERT INTO reviews (reviewer_id, solution_id, content) VALUES (?, ?, ?)",
			review.ReviewerID, review.SolutionID, review.Content)

		if err != nil {
			return result, err
		}
		insertedID, err := res.LastInsertId()
		if err != nil {
			return result, err
		}
		result.ID = strconv.FormatInt(insertedID, 10)
	case "postgres":
		query := "INSERT INTO reviews (reviewer_id, solution_id, content) VALUES ($1, $2, $3) RETURNING id"
		stmt, err := r.db.Prepare(query)
		if err != nil {
			return result, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(review.ReviewerID, review.SolutionID, review.Content).Scan(&result.ID)
		if err != nil {
			return result, err
		}

		return result, nil
	default:
		return result, fmt.Errorf("Unknown dbtype")

	}

	return result, nil
}

func (r *ReviewRepository) Get(ID string) (mathbattle.Review, error) {
	row := r.db.QueryRow("SELECT id, reviewer_id, solution_id, content FROM reviews WHERE id = ?", ID)
	result := mathbattle.Review{}
	err := row.Scan(&result.ID, &result.ReviewerID, &result.SolutionID, &result.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Review{}, mathbattle.ErrNotFound
		}
		return mathbattle.Review{}, err
	}
	result.ID = ID

	return result, nil
}

func (r *ReviewRepository) FindMany(reviewerID, solutionID string) ([]mathbattle.Review, error) {
	query := "SELECT id, reviewer_id, solution_id, content FROM reviews"
	query, whereArgs := createWhereClause(query, []whereDescriptor{
		{"reviewer_id", reviewerID},
		{"solution_id", solutionID}})

	rows, err := r.db.Query(query, whereArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []mathbattle.Review{}, mathbattle.ErrNotFound
		}
	}
	defer rows.Close()

	result := []mathbattle.Review{}

	for rows.Next() {
		curReview := mathbattle.Review{}

		err = rows.Scan(&curReview.ID, &curReview.ReviewerID, &curReview.SolutionID, &curReview.Content)
		if err != nil {
			if err == sql.ErrNoRows {
				return []mathbattle.Review{}, nil
			}
			return result, err
		}

		result = append(result, curReview)
	}

	return result, nil
}

func (r *ReviewRepository) Update(review mathbattle.Review) error {
	_, err := r.db.Exec("UPDATE reviews SET reviewer_id = $1, solution_id = $2, content = $3 WHERE id = $4",
		review.ReviewerID, review.SolutionID, review.Content, review.ID)
	return err
}

func (r *ReviewRepository) Delete(ID string) error {
	_, err := r.db.Exec("DELETE FROM reviews WHERE id = $1", ID)
	return err
}
