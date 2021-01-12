package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var gDB *sql.DB = nil

func initSqliteDb(dbPath string) error {
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

func initPostgresDb(connectionString string) error {
	if gDB == nil {
		dbName, err := getDbNameFromConnString(connectionString)
		if err != nil {
			return err
		}

		genericConnString := removeDbNameFromConnString(connectionString)

		db, err := sql.Open("postgres", genericConnString)
		if err != nil {
			log.Printf("Init, sql.Open() error: %v", err)
			return err
		}

		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			pgerr, ok := err.(*pq.Error)
			if !ok {
				return err
			}

			if pgerr.Code != "42P04" { // Duplicate database
				log.Printf("Init, failed to create database error: %v", err)
				return err
			}

			log.Printf("Don't need to create database, already exists")
		}

		err = db.Close()
		if err != nil {
			return err
		}

		gDB, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Printf("Init, sql.Open() error: %v", err)
			return err
		}
	}

	return nil
}

func getDbNameFromConnString(connectionString string) (string, error) {
	for _, part := range strings.Split(connectionString, " ") {
		if strings.HasPrefix(part, "dbname=") {
			dbnameParts := strings.Split(part, "=")
			if len(dbnameParts) != 2 {
				return "", fmt.Errorf("Unexpected connection string format")
			}

			return dbnameParts[1], nil
		}
	}

	return "", fmt.Errorf("Unexpected connection string format")
}

func removeDbNameFromConnString(connectionString string) string {
	newParts := []string{}
	for _, part := range strings.Split(connectionString, " ") {
		if strings.HasPrefix(part, "dbname=") {
			continue
		}
		newParts = append(newParts, part)
	}

	return strings.Join(newParts, " ")
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

func DeinitAndRemove(dbType, connectionString string) error {
	if gDB != nil {
		err := gDB.Close()
		if err != nil {
			return err
		}
		gDB = nil
	}

	switch dbType {
	case "sqlite3":
		return os.RemoveAll(connectionString)
	case "postgres":
		genericConnString := removeDbNameFromConnString(connectionString)
		db, err := sql.Open("postgres", genericConnString)
		if err != nil {
			log.Printf("Init, sql.Open() error: %v", err)
			return err
		}

		dbName, err := getDbNameFromConnString(connectionString)
		if err != nil {
			return err
		}

		log.Printf("Removing database %v", dbName)
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))

		return err
	default:
		return fmt.Errorf("Unknown database type")
	}
}

type sqlRepository struct {
	db     *sql.DB
	dbType string
}

func newSqlRepository(dbType, connectionString string) (sqlRepository, error) {
	if gDB == nil {
		switch dbType {
		case "sqlite3":
			if err := initSqliteDb(connectionString); err != nil {
				return sqlRepository{}, err
			}
		case "postgres":
			if err := initPostgresDb(connectionString); err != nil {
				return sqlRepository{}, err
			}
		default:
			return sqlRepository{}, fmt.Errorf("Unknown repository type")
		}
	}

	return sqlRepository{
		db:     gDB,
		dbType: dbType,
	}, nil
}

type whereDescriptor struct {
	ParamName  string
	ParamValue string
}

func createWhereClause(query string, whereDescriptors []whereDescriptor) (string, []interface{}) {
	i := 1
	whereClauses := []string{}
	whereArgs := []interface{}{}
	for _, whereArg := range whereDescriptors {
		if whereArg.ParamValue != "" {
			whereClauses = append(whereClauses, fmt.Sprintf(" %s = $%d", whereArg.ParamName, i))
			whereArgs = append(whereArgs, whereArg.ParamValue)
			i++
		}
	}

	if len(whereClauses) != 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, whereArgs
}
