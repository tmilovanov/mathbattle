package sqlite

import (
	"database/sql"

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
