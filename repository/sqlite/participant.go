package sqlite

import (
	"database/sql"
	"errors"
	"strconv"

	mathbattle "mathbattle/models"
)

type ParticipantRepository struct {
	sqliteRepository
}

func NewParticipantRepository(dbPath string) (ParticipantRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return ParticipantRepository{}, err
	}

	return ParticipantRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func NewParticipantRepositoryTemp(dbName string) (ParticipantRepository, error) {
	sqliteRepository, err := newTempSqliteRepository(dbName)
	if err != nil {
		return ParticipantRepository{}, err
	}

	return ParticipantRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func (r *ParticipantRepository) Store(participant mathbattle.Participant) (mathbattle.Participant, error) {
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

func (r *ParticipantRepository) GetByID(ID string) (mathbattle.Participant, error) {
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

func (r *ParticipantRepository) GetByTelegramID(telegramID int64) (mathbattle.Participant, error) {
	result := mathbattle.Participant{}

	row := r.db.QueryRow("SELECT id, tg_chat_id, name, school, grade, register_time FROM participants WHERE tg_chat_id = ?", telegramID)
	err := row.Scan(&result.ID, &result.TelegramID, &result.Name, &result.School, &result.Grade, &result.RegistrationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mathbattle.Participant{}, mathbattle.ErrNotFound
		}
		return mathbattle.Participant{}, err
	}

	return result, nil
}

func (r *ParticipantRepository) GetAll() ([]mathbattle.Participant, error) {
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

func (r *ParticipantRepository) Update(participant mathbattle.Participant) error {
	_, err := r.db.Exec("UPDATE participants SET tg_chat_id = ?, name = ?, grade = ?, school = ?, register_time = ? WHERE id = ?",
		participant.TelegramID, participant.Name, participant.Grade, participant.School, participant.RegistrationTime,
		participant.ID)
	return err
}

func (r *ParticipantRepository) Delete(ID string) error {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM participants WHERE id = ?", intID)
	return err
}
