package sqldb

import (
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
			content TEXT,
			juri_comment TEXT,
			mark INTEGER
		)`
	case "postgres":
		createStmt = `CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL UNIQUE,
			reviewer_id INTEGER,
			solution_id INTEGER,
			content TEXT,
			juri_comment TEXT,
			mark INTEGER
		)`
	}
	_, err := r.db.Exec(createStmt)
	return err
}

func (r *ReviewRepository) Store(review mathbattle.Review) (mathbattle.Review, error) {
	result := review

	switch r.dbType {
	case "sqlite3":
		res, err := r.db.Exec("INSERT INTO reviews (reviewer_id, solution_id, content, juri_comment, mark) VALUES (?, ?, ?, ?, ?)",
			review.ReviewerID, review.SolutionID, review.Content, review.JuriComment, review.Mark)

		if err != nil {
			return result, err
		}
		insertedID, err := res.LastInsertId()
		if err != nil {
			return result, err
		}
		result.ID = strconv.FormatInt(insertedID, 10)
	case "postgres":
		query := "INSERT INTO reviews (reviewer_id, solution_id, content, juri_comment, mark) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		stmt, err := r.db.Prepare(query)
		if err != nil {
			return result, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(review.ReviewerID, review.SolutionID, review.Content, review.JuriComment, review.Mark).Scan(&result.ID)
		if err != nil {
			return result, err
		}

		return result, nil
	default:
		return result, fmt.Errorf("Unknown dbtype")

	}

	return result, nil
}

func (r *ReviewRepository) getManyWhere(whereStr string, whereArgs ...interface{}) ([]mathbattle.Review, error) {
	result := []mathbattle.Review{}

	rows, err := r.db.Query("SELECT id, reviewer_id, solution_id, content, juri_comment, mark FROM reviews WHERE "+whereStr, whereArgs...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var cur mathbattle.Review
		err = rows.Scan(&cur.ID, &cur.ReviewerID, &cur.SolutionID, &cur.Content, &cur.JuriComment, &cur.Mark)
		if err != nil {
			return result, err
		}

		result = append(result, cur)
	}

	return result, nil
}

func (r *ReviewRepository) getOneWhere(whereStr string, whereArgs ...interface{}) (mathbattle.Review, error) {
	res, err := r.getManyWhere(whereStr, whereArgs...)
	if err != nil {
		return mathbattle.Review{}, err
	}

	if len(res) == 0 {
		return mathbattle.Review{}, mathbattle.ErrNotFound
	}

	return res[0], nil
}

func (r *ReviewRepository) Get(ID string) (mathbattle.Review, error) {
	return r.getOneWhere("id = $1", ID)
}

func (r *ReviewRepository) FindMany(reviewerID, solutionID string) ([]mathbattle.Review, error) {
	whereClause, whereArgs := joinWhereOmitEmpty([]whereDescriptor{
		{"reviewer_id", reviewerID},
		{"solution_id", solutionID}})
	return r.getManyWhere(whereClause, whereArgs...)
}

func (r *ReviewRepository) Update(review mathbattle.Review) error {
	_, err := r.db.Exec(`
	UPDATE reviews SET
		reviewer_id = $1, solution_id = $2, content = $3, juri_comment = $4, mark = $5
	WHERE 
		id = $6`,
		review.ReviewerID, review.SolutionID, review.Content, review.JuriComment, review.Mark, review.ID)
	return err
}

func (r *ReviewRepository) Delete(ID string) error {
	_, err := r.db.Exec("DELETE FROM reviews WHERE id = $1", ID)
	return err
}
