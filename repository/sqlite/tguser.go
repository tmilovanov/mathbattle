package sqlite

import (
	"database/sql"
	"errors"
	mathbattle "mathbattle/models"
	"strconv"
)

type TelegramUserRepository struct {
	sqliteRepository
}

func NewTelegramUserRepository(dbPath string) (TelegramUserRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return TelegramUserRepository{}, err
	}

	return TelegramUserRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func (r *TelegramUserRepository) Store(user mathbattle.TelegramUser) (mathbattle.TelegramUser, error) {
	result := user

	res, err := r.db.Exec("INSERT INTO tgusers (tg_chat_id, is_admin) VALUES (?, ?)",
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

func (r *TelegramUserRepository) GetByID(ID string) (mathbattle.TelegramUser, error) {
	result := mathbattle.TelegramUser{}

	intID, err := strconv.Atoi(ID)
	if err != nil {
		return result, err
	}

	row := r.db.QueryRow("SELECT id, tg_chat_id, is_admin FROM tgusers WHERE id = ?", intID)
	err = row.Scan(&result.ID, &result.ChatID, &result.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

	return result, nil
}

func (r *TelegramUserRepository) GetByTelegramID(ID int64) (mathbattle.TelegramUser, error) {
	result := mathbattle.TelegramUser{}

	row := r.db.QueryRow("SELECT id, tg_chat_id, is_admin FROM tgusers WHERE tg_chat_id = ?", ID)
	err := row.Scan(&result.ID, &result.ChatID, &result.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

	return result, nil
}

func (r *TelegramUserRepository) GetOrCreateByTelegramID(ID int64) (mathbattle.TelegramUser, error) {
	user, err := r.GetByTelegramID(ID)
	if err == nil {
		return user, nil
	}

	if err == mathbattle.ErrNotFound {
		return r.Store(mathbattle.TelegramUser{ChatID: ID})
	}

	return user, err
}
