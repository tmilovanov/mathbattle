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
	if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Dir(dbPath), 0777); err != nil {
			return sqliteRepository{}, err
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return sqliteRepository{}, err
	}

	return sqliteRepository{
		db: db,
	}, nil
}
