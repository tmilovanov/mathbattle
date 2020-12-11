package sqlite

import (
	"database/sql"
	"errors"
	mathbattle "mathbattle/models"
	"strconv"
	"strings"
)

type ScheduledMessageRepository struct {
	sqliteRepository
}

func NewScheduledMessageRepository(dbPath string) (ScheduledMessageRepository, error) {
	sqliteRepository, err := newSqliteRepository(dbPath)
	if err != nil {
		return ScheduledMessageRepository{}, err
	}

	return ScheduledMessageRepository{
		sqliteRepository: sqliteRepository,
	}, nil
}

func serializeRecievers(recievers []string) string {
	return strings.Join(recievers, ",")
}

func deserializeRecievers(recievers string) []string {
	return strings.Split(recievers, ",")
}

func (r *ScheduledMessageRepository) Store(msg mathbattle.ScheduledMessage) (mathbattle.ScheduledMessage, error) {
	result := msg
	res, err := r.db.Exec("INSERT INTO scheduled_messages (type, send_time, content, recievers) VALUES (?, ?, ?, ?)",
		msg.RecieversType, msg.SendTime, msg.Message, serializeRecievers(msg.Recievers))

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

func (r *ScheduledMessageRepository) Get(ID string) (mathbattle.ScheduledMessage, error) {
	result := mathbattle.ScheduledMessage{}

	row := r.db.QueryRow("SELECT type, send_time, content, recievers FROM scheduled_messages WHERE id = ?", ID)
	var recievers string
	err := row.Scan(&result.RecieversType, &result.SendTime, &result.Message, &recievers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, mathbattle.ErrNotFound
		}
		return result, err
	}
	result.ID = ID
	result.Recievers = deserializeRecievers(recievers)

	return result, nil
}

func (r *ScheduledMessageRepository) GetAll() ([]mathbattle.ScheduledMessage, error) {
	result := []mathbattle.ScheduledMessage{}
	rows, err := r.db.Query("SELECT id, type, send_time, content, recievers FROM scheduled_messages")
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var recievers string
		curParticipant := mathbattle.ScheduledMessage{}
		err = rows.Scan(&curParticipant.ID, &curParticipant.RecieversType, &curParticipant.SendTime, &curParticipant.Message,
			&recievers)
		if err != nil {
			return result, err
		}
		curParticipant.Recievers = deserializeRecievers(recievers)

		result = append(result, curParticipant)
	}

	return result, nil
}

func (r *ScheduledMessageRepository) Update(msg mathbattle.ScheduledMessage) error {
	_, err := r.db.Exec("UPDATE scheduled_messages SET type = ?, send_time = ?, content = ?, recievers = ? WHRE id = ?",
		msg.RecieversType, msg.SendTime, msg.Message, serializeRecievers(msg.Recievers), msg.ID)
	return err
}

func (r *ScheduledMessageRepository) Delete(ID string) error {
	_, err := r.db.Exec("DELETE FROM scheduled_messages WHERE id = ?", ID)
	return err
}
