package sqlite

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	mathbattle "mathbattle/models"
)

type ReviewRepository struct {
	sqliteRepository
}

func NewReviewRepository(dbPath string) (ReviewRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return ReviewRepository{}, err
	}

	return ReviewRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func NewReviewRepositoryTemp(dbName string) (ReviewRepository, error) {
	sqliteRepository, err := newTempSqliteRepository(dbName)
	if err != nil {
		return ReviewRepository{}, err
	}

	return ReviewRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func (r *ReviewRepository) Store(review mathbattle.Review) (mathbattle.Review, error) {
	result := review
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
	whereClauses := []string{}
	whereArgs := []interface{}{}
	if reviewerID != "" {
		whereClauses = append(whereClauses, " reviewer_id = ?")
		whereArgs = append(whereArgs, reviewerID)
	}
	if solutionID != "" {
		whereClauses = append(whereClauses, " solution_id = ?")
		whereArgs = append(whereArgs, solutionID)
	}
	if len(whereClauses) != 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

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
				return result, mathbattle.ErrNotFound
			}
			return result, err
		}

		result = append(result, curReview)
	}

	return result, nil
}

func (r *ReviewRepository) Update(review mathbattle.Review) error {
	_, err := r.db.Exec("UPDATE reviews SET reviewer_id = ?, solution_id = ?, content = ? WHERE id = ?",
		review.ReviewerID, review.SolutionID, review.Content, review.ID)
	return err
}

func (r *ReviewRepository) Delete(ID string) error {
	_, err := r.db.Exec("DELETE FROM reviews WHERE id = ?", ID)
	return err
}
