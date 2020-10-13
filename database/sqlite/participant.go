package sqlite

import (
	"database/sql"
	"errors"
	"strconv"

	mathbattle "mathbattle/models"
)

type SQLParticipantRepository struct {
	sqliteRepository
}

func NewSQLParticipantRepository(dbPath string) (SQLParticipantRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return SQLParticipantRepository{}, err
	}

	return SQLParticipantRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func NewSQLParticipantRepositoryTemp(dbName string) (SQLParticipantRepository, error) {
	sqliteRepository, err := newTempSqliteRepository(dbName)
	if err != nil {
		return SQLParticipantRepository{}, err
	}

	return SQLParticipantRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func (r *SQLParticipantRepository) Store(participant mathbattle.Participant) (mathbattle.Participant, error) {
	result := participant
	res, err := r.db.Exec("INSERT INTO participants (tg_chat_id, name, school, grade, register_time) VALUES (?, ?, ?, ?, ?)",
		participant.TelegramID, participant.Name, participant.School, participant.Grade, participant.RegistrationTime)

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

func (r *SQLParticipantRepository) GetByID(ID string) (mathbattle.Participant, error) {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return mathbattle.Participant{}, err
	}

	row := r.db.QueryRow("SELECT tg_chat_id, name, school, grade, register_time FROM participants WHERE id = ?", intID)
	result := mathbattle.Participant{}
	err = row.Scan(&result.TelegramID, &result.Name, &result.School, &result.Grade, &result.RegistrationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Participant{}, mathbattle.ErrNotFound
		}
		return mathbattle.Participant{}, err
	}
	result.ID = ID

	return result, nil
}

func (r *SQLParticipantRepository) GetByTelegramID(telegramID string) (mathbattle.Participant, error) {
	row := r.db.QueryRow("SELECT id, name, school, grade, register_time FROM participants WHERE tg_chat_id = ?", telegramID)
	var id int
	result := mathbattle.Participant{}
	err := row.Scan(&id, &result.Name, &result.School, &result.Grade, &result.RegistrationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Participant{}, mathbattle.ErrNotFound
		}
		return mathbattle.Participant{}, err
	}
	result.ID = strconv.Itoa(id)

	return result, nil
}

func (r *SQLParticipantRepository) GetAll() ([]mathbattle.Participant, error) {
	rows, err := r.db.Query("SELECT id, tg_chat_id, name, school, grade, register_time FROM participants")
	if err != nil {
		return []mathbattle.Participant{}, err
	}
	defer rows.Close()

	result := []mathbattle.Participant{}
	for rows.Next() {
		var id int
		curParticipant := mathbattle.Participant{}
		err = rows.Scan(&id, &curParticipant.TelegramID, &curParticipant.Name, &curParticipant.School,
			&curParticipant.Grade, &curParticipant.Name)
		if err != nil {
			return []mathbattle.Participant{}, err
		}
		curParticipant.ID = strconv.Itoa(id)
		result = append(result, curParticipant)
	}
	return result, nil
}

func (r *SQLParticipantRepository) Delete(ID string) error {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM participants WHERE id = ?", intID)
	return err
}
