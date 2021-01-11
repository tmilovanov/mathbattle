package sqlite

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var gDB *sql.DB = nil

func Init(dbPath string) error {
	log.Printf("Init db: %v", dbPath)
	if gDB == nil {
		if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
			log.Printf("%v not exist, creating.", filepath.Dir(dbPath))
			if err := os.MkdirAll(filepath.Dir(dbPath), 0766); err != nil {
				log.Printf("Fld to create %v, error: %v", filepath.Dir(dbPath), err)
				return err
			}
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Printf("Init, sql.Open() error: %v", err)
			return err
		}

		gDB = db
	}

	return nil
}

func Deinit() error {
	log.Printf("Deinit db")
	if gDB != nil {
		err := gDB.Close()
		gDB = nil
		return err
	}

	return nil
}

type sqliteRepository struct {
	db *sql.DB
}

func newSqliteRepository(dbPath string) (sqliteRepository, error) {
	if gDB == nil {
		if err := Init(dbPath); err != nil {
			return sqliteRepository{}, err
		}
	}

	return sqliteRepository{
		db: gDB,
	}, nil
}
