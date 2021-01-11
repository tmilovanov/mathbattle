package sqlite

import (
	"database/sql"
	"errors"
	"strconv"

	"mathbattle/models/mathbattle"
)

type UserRepository struct {
	sqliteRepository
	participantRepository *ParticipantRepository
}

func NewUserRepository(dbPath string) (*UserRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return nil, err
	}

	result := &UserRepository{
		sqliteRepository: sqliteRepository,
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
	createStmt := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			tg_chat_id VARCHAR(64) UNIQUE,
			tg_name VARCHAR(100) UNIQUE,
			is_admin BOOL,
			registration_time DATETIME
		)`
	_, err := r.db.Exec(createStmt)
	return err
}

func (r *UserRepository) Store(user mathbattle.User) (mathbattle.User, error) {
	result := user

	res, err := r.db.Exec("INSERT INTO users (tg_chat_id, tg_name, is_admin, registration_time) VALUES (?, ?, ?, ?)",
		user.TelegramID, user.TelegramName, user.IsAdmin, user.RegistrationTime)
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

func (r *UserRepository) getWhere(whereStr string, whereArgs ...interface{}) (mathbattle.User, error) {
	result := mathbattle.User{}
	row := r.db.QueryRow("SELECT id, tg_chat_id, tg_name, is_admin, registration_time FROM users WHERE "+whereStr, whereArgs...)
	err := row.Scan(&result.ID, &result.TelegramID, &result.TelegramName, &result.IsAdmin, &result.RegistrationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

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
	return r.getWhere("id = ?", ID)
}

func (r *UserRepository) GetByTelegramID(ID int64) (mathbattle.User, error) {
	return r.getWhere("tg_chat_id = ?", ID)
}

func (r *UserRepository) GetByTelegramName(name string) (mathbattle.User, error) {
	return r.getWhere("tg_name = ?", name)
}

func (r *UserRepository) Update(user mathbattle.User) error {
	_, err := r.db.Exec("UPDATE users SET tg_chat_id = ?, tg_name = ?, is_admin = ?, registration_time = ? WHERE id = ?",
		user.TelegramID, user.TelegramName, user.IsAdmin, user.RegistrationTime, user.ID)
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

	_, err = r.db.Exec("DELETE FROM users WHERE id = ?", user.ID)
	return err
}
