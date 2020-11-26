package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteRepository struct {
	db *sql.DB
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

func newTempSqliteRepository(dbName string) (sqliteRepository, error) {
	dbPath := filepath.Join(os.TempDir(), dbName)
	result, err := newSqliteRepository(dbPath)
	if err != nil {
		return result, err
	}

	err = result.deleteAllTables()
	if err != nil {
		return result, err
	}

	return result, result.CreateFirstTime()
}

func (r *sqliteRepository) deleteAllTables() error {
	tableDeleters := []string{
		"DROP TABLE tgusers",
		"DROP TABLE participants",
		"DROP TABLE problems",
		"DROP TABLE rounds",
		"DROP TABLE solutions",
		"DROP TABLE reviews",
	}
	for _, deleteStmt := range tableDeleters {
		_, _ = r.db.Exec(deleteStmt)
	}
	return nil
}

func (r *sqliteRepository) CreateFirstTime() error {
	tableCreators := []string{
		`CREATE TABLE IF NOT EXISTS tgusers (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			tg_chat_id VARCHAR(100),
			is_admin BOOL
		)`,
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
			solve_start DATETIME,
			solve_end DATETIME,
			review_start DATETIME,
			review_end DATETIME,
			problems_distribution TEXT,
			solutions_distribution TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS solutions (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			round_id INTEGER,
			participant_id INTEGER,
			problem_id INTEGER,
			parts TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			reviewer_id INTEGER,
			solution_id INTEGER,
			content TEXT
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
