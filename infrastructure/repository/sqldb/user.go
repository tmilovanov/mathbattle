package sqldb

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"mathbattle/models/mathbattle"
)

type UserRepository struct {
	sqlRepository
	participantRepository *ParticipantRepository
}

func NewUserRepository(dbType, connectionString string) (*UserRepository, error) {
	sqlRepository, err := newSqlRepository(dbType, connectionString)
	if err != nil {
		return nil, err
	}

	result := &UserRepository{
		sqlRepository: sqlRepository,
	}

	if err := result.CreateTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *UserRepository) SetParticipantRepository(pr *ParticipantRepository) {
	r.participantRepository = pr
}

func (r *UserRepository) CreateTable() error {
	var createStmt string

	switch r.dbType {
	case "sqlite3":
		createStmt = `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			tg_chat_id VARCHAR(64) UNIQUE,
			tg_firstname VARCHAR(100),
			tg_lastname VARCHAR(100),
			tg_username VARCHAR(100),
			is_admin BOOL,
			registration_time DATETIME
		)`
	case "postgres":
		createStmt = `CREATE TABLE IF NOT EXISTS users (
			id SERIAL UNIQUE,
			tg_chat_id VARCHAR(64) UNIQUE,
			tg_firstname VARCHAR(100),
			tg_lastname VARCHAR(100),
			tg_username VARCHAR(100),
			is_admin BOOL,
			registration_time TIMESTAMP
		)`
	}

	_, err := r.db.Exec(createStmt)
	return err
}

func (r *UserRepository) Store(user mathbattle.User) (mathbattle.User, error) {
	result := user

	switch r.dbType {
	case "sqlite3":
		res, err := r.db.Exec(`
			INSERT INTO users (tg_chat_id, tg_firstname, tg_lastname, tg_username, is_admin, registration_time)
			VALUES (?, ?, ?, ?, ?, ?)`,
			user.TelegramID, user.TelegramFirstName, user.TelegramLastName, user.TelegramUsername, user.IsAdmin, user.RegistrationTime)
		if err != nil {
			return result, err
		}

		insertedID, err := res.LastInsertId()
		if err != nil {
			return result, err
		}
		result.ID = strconv.FormatInt(insertedID, 10)

		return result, nil
	case "postgres":
		query := `INSERT INTO users (tg_chat_id, tg_firstname, tg_lastname, tg_username, is_admin, registration_time)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		stmt, err := r.db.Prepare(query)
		if err != nil {
			return result, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(user.TelegramID, user.TelegramFirstName, user.TelegramLastName, user.TelegramUsername,
			user.IsAdmin, user.RegistrationTime).Scan(&result.ID)
		if err != nil {
			return result, err
		}

		return result, nil
	default:
		return result, fmt.Errorf("Unexpected db type")
	}
}

func (r *UserRepository) getWhere(whereStr string, whereArgs ...interface{}) (mathbattle.User, error) {
	result := mathbattle.User{}
	row := r.db.QueryRow(`SELECT
		id, tg_chat_id, tg_firstname, tg_lastname, tg_username, is_admin, registration_time
	FROM users WHERE `+whereStr, whereArgs...)
	err := row.Scan(&result.ID, &result.TelegramID, &result.TelegramFirstName, &result.TelegramLastName, &result.TelegramUsername,
		&result.IsAdmin, &result.RegistrationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}
	result.SetRegistrationTime(result.RegistrationTime)

	return result, nil
}

func (r *UserRepository) GetAll() ([]mathbattle.User, error) {
	result := []mathbattle.User{}
	rows, err := r.db.Query("SELECT id FROM users")
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

		cur, err := r.GetByID(curID)
		if err != nil {
			return result, err
		}

		result = append(result, cur)
	}

	return result, nil
}

func (r *UserRepository) GetByID(ID string) (mathbattle.User, error) {
	return r.getWhere("id = $1", ID)
}

func (r *UserRepository) GetByTelegramID(ID int64) (mathbattle.User, error) {
	return r.getWhere("tg_chat_id = $1", ID)
}

func (r *UserRepository) GetByTelegramName(name string) (mathbattle.User, error) {
	return r.getWhere("tg_name = $1", name)
}

func (r *UserRepository) Update(user mathbattle.User) error {
	_, err := r.db.Exec(`UPDATE users SET 
		tg_chat_id = $1, tg_firstname=$2, tg_lastname=$3, tg_username = $4, is_admin = $5, registration_time = $6 WHERE id = $7`,
		user.TelegramID, user.TelegramFirstName, user.TelegramLastName, user.TelegramUsername,
		user.IsAdmin, user.RegistrationTime, user.ID)
	return err
}

func (r *UserRepository) Delete(user mathbattle.User) error {
	participant, err := r.participantRepository.GetByUserID(user.ID)
	if err != nil {
		return err
	}

	err = r.participantRepository.Delete(participant.ID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM users WHERE id = $1", user.ID)
	return err
}
