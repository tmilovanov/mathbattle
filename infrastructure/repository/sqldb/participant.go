package sqldb

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"mathbattle/models/mathbattle"
)

type ParticipantRepository struct {
	sqlRepository
	userRepository *UserRepository
}

func NewParticipantRepository(dbType, connectionString string, userRepository *UserRepository) (*ParticipantRepository, error) {
	sqlRepository, err := newSqlRepository(dbType, connectionString)
	if err != nil {
		return nil, err
	}

	result := &ParticipantRepository{
		sqlRepository:  sqlRepository,
		userRepository: userRepository,
	}

	if err := result.CreateTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ParticipantRepository) CreateTable() error {
	var createStmt string
	switch r.dbType {
	case "sqlite3":
		createStmt = `CREATE TABLE IF NOT EXISTS participants (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			user_id INTEGER NOT NULL,
			name VARCHAR(100),
			school VARCHAR(256),
			grade INTEGER,
			is_active BOOL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`
	case "postgres":
		createStmt = `CREATE TABLE IF NOT EXISTS participants (
			id SERIAL UNIQUE,
			user_id INTEGER NOT NULL,
			name VARCHAR(100),
			school VARCHAR(256),
			grade INTEGER,
			is_active BOOL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`
	}

	_, err := r.db.Exec(createStmt)

	return err
}

func (r *ParticipantRepository) Store(participant mathbattle.Participant) (mathbattle.Participant, error) {
	result := participant

	switch r.dbType {
	case "sqlite3":
		res, err := r.db.Exec("INSERT INTO participants (user_id, name, school, grade, is_active) VALUES (?, ?, ?, ?, ?)",
			participant.User.ID, participant.Name, participant.School, participant.Grade, participant.IsActive)

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
		query := "INSERT INTO participants (user_id, name, school, grade, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		stmt, err := r.db.Prepare(query)
		if err != nil {
			return result, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(participant.User.ID, participant.Name, participant.School, participant.Grade, participant.IsActive).Scan(&result.ID)
		if err != nil {
			return result, err
		}

		return result, nil
	default:
		return result, fmt.Errorf("Unexpected db type")
	}
}

func (r *ParticipantRepository) getWhere(whereStr string, whereArgs ...interface{}) (mathbattle.Participant, error) {
	result := mathbattle.Participant{}
	row := r.db.QueryRow("SELECT id, user_id, name, school, grade, is_active FROM participants WHERE "+whereStr, whereArgs...)
	err := row.Scan(&result.ID, &result.User.ID, &result.Name, &result.School, &result.Grade, &result.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}

	user, err := r.userRepository.GetByID(result.User.ID)
	if err != nil {
		return result, err
	}

	result.User = user

	return result, nil
}

func (r *ParticipantRepository) GetByID(ID string) (mathbattle.Participant, error) {
	return r.getWhere("id = $1", ID)
}

func (r *ParticipantRepository) GetByTelegramID(telegramID int64) (mathbattle.Participant, error) {
	user, err := r.userRepository.GetByTelegramID(telegramID)
	if err != nil {
		return mathbattle.Participant{}, err
	}

	return r.getWhere("user_id = $1", user.ID)
}

func (r *ParticipantRepository) GetByUserID(userID string) (mathbattle.Participant, error) {
	return r.getWhere("user_id = $1", userID)
}

func (r *ParticipantRepository) GetAll() ([]mathbattle.Participant, error) {
	rows, err := r.db.Query("SELECT id, user_id, name, school, grade, is_active FROM participants")
	if err != nil {
		return []mathbattle.Participant{}, err
	}
	defer rows.Close()

	result := []mathbattle.Participant{}
	for rows.Next() {
		curParticipant := mathbattle.Participant{}
		err = rows.Scan(&curParticipant.ID, &curParticipant.User.ID, &curParticipant.Name, &curParticipant.School,
			&curParticipant.Grade, &curParticipant.IsActive)
		if err != nil {
			return []mathbattle.Participant{}, err
		}

		user, err := r.userRepository.GetByID(curParticipant.User.ID)
		if err != nil {
			return []mathbattle.Participant{}, err
		}

		curParticipant.User = user

		result = append(result, curParticipant)
	}
	return result, nil
}

func (r *ParticipantRepository) Update(participant mathbattle.Participant) error {
	_, err := r.db.Exec("UPDATE participants SET user_id = $1, name = $2, grade = $3, school = $4, is_active = $5 WHERE id = $6",
		participant.User.ID, participant.Name, participant.Grade, participant.School, participant.IsActive,
		participant.ID)
	if err != nil {
		return err
	}

	return r.userRepository.Update(participant.User)
}

func (r *ParticipantRepository) Delete(ID string) error {
	_, err := r.db.Exec("DELETE FROM participants WHERE id = $1", ID)
	return err
}
