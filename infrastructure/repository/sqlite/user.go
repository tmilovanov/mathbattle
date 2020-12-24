package sqlite

import (
	"database/sql"
	"errors"
	"strconv"

	"mathbattle/models/mathbattle"
)

type UserRepository struct {
	sqliteRepository
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

func (r *UserRepository) CreateTable() error {
	createStmt := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			tg_name VARCHAR(100),
			tg_chat_id VARCHAR(100),
			is_admin BOOL
		)`
	_, err := r.db.Exec(createStmt)
	return err
}

func (r *UserRepository) Store(user mathbattle.User) (mathbattle.User, error) {
	result := user

	res, err := r.db.Exec("INSERT INTO users (tg_chat_id, is_admin) VALUES (?, ?)",
		user.ChatID, false)
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

func (r *UserRepository) GetByID(ID string) (mathbattle.User, error) {
	result := mathbattle.User{}

	intID, err := strconv.Atoi(ID)
	if err != nil {
		return result, err
	}

	row := r.db.QueryRow("SELECT id, tg_chat_id, is_admin FROM users WHERE id = ?", intID)
	err = row.Scan(&result.ID, &result.ChatID, &result.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

	return result, nil
}

func (r *UserRepository) GetByTelegramID(ID int64) (mathbattle.User, error) {
	result := mathbattle.User{}

	row := r.db.QueryRow("SELECT id, tg_chat_id, is_admin FROM users WHERE tg_chat_id = ?", ID)
	err := row.Scan(&result.ID, &result.ChatID, &result.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

	return result, nil
}

func (r *UserRepository) GetOrCreateByTelegramID(ID int64) (mathbattle.User, error) {
	user, err := r.GetByTelegramID(ID)
	if err == nil {
		return user, nil
	}

	if err == mathbattle.ErrNotFound {
		return r.Store(mathbattle.User{ChatID: ID})
	}

	return user, err
}
