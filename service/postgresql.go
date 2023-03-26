package service

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/paulsmith/gogeos/geos"
	"github.com/twinj/uuid"
	"log"
	"time"
)

const (
	insertMessage  = "INSERT INTO message (id, user_id, content, location, client_id, sent_at, received_at) VALUES ($1, $2, $3, ST_GeomFromText($4), $5, $6, $7) RETURNING created_at"
	selectMessage  = "SELECT m.id, l.id as userId, l.username, m.content, m.location, m.created_at, m.client_id, m.sent_at, m.received_at FROM message m JOIN login l on m.user_id = l.id WHERE m.id = $1"
	selectMessages = "SELECT m.id, l.id as userId, l.username, m.content, m.location, m.created_at, m.client_id, m.sent_at, m.received_at FROM message m JOIN login l on m.user_id = l.id WHERE ST_DistanceSphere(m.location, $1) <= $2 AND m.created_at > $3 ORDER BY m.created_at LIMIT $4"
)

type postgresqlMessageRepository struct {
	db *sql.DB
}

func (p *postgresqlMessageRepository) AddMessage(message Message) (StoredMessage, error) {
	id := uuid.NewV4().String()

	receivedAt := time.Now().UTC()
	row := p.db.QueryRow(insertMessage, id, message.Sender.Id, message.Content, fmt.Sprintf("POINT (%f %f)",
		message.Location.Long, message.Location.Lat), message.ClientId, message.SentAt, receivedAt)

	var createdAt time.Time
	err := row.Scan(&createdAt)

	if err != nil {
		return StoredMessage{}, err
	}

	return StoredMessage{id, createdAt, receivedAt, message}, nil
}

func (p *postgresqlMessageRepository) GetMessage(id string) (StoredMessage, error) {
	row := p.db.QueryRow(selectMessage, id)

	loc := make([]byte, 0)
	var message StoredMessage
	err := row.Scan(&message.Id, &message.Sender.Id, &message.Sender.Username, &message.Content, &loc,
		&message.CreatedAt, &message.ClientId, &message.SentAt, &message.ReceivedAt)

	if err == sql.ErrNoRows {
		return StoredMessage{}, nil
	} else if err != nil {
		return StoredMessage{}, newErrRepository(err.Error())
	}

	geom, err := geos.FromHex(string(loc))

	if err != nil {
		return StoredMessage{}, newErrRepository(err.Error())
	}

	message.Location.Long, err = geom.X()

	if err != nil {
		return StoredMessage{}, newErrRepository(err.Error())
	}

	message.Location.Lat, err = geom.Y()

	if err != nil {
		return StoredMessage{}, newErrRepository(err.Error())
	}

	return message, nil
}

func (p *postgresqlMessageRepository) GetMessagesForLocation(location Location, radiusMeters float64, limit int, after time.Time) ([]StoredMessage, error) {
	rows, err := p.db.Query(selectMessages, fmt.Sprintf("POINT (%f %f)", location.Long, location.Lat),
		radiusMeters, after, limit)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	messages := make([]StoredMessage, 0)

	for rows.Next() {
		loc := make([]byte, 0)
		var message StoredMessage
		err := rows.Scan(&message.Id, &message.Sender.Id, &message.Sender.Username, &message.Content, &loc,
			&message.CreatedAt, &message.ClientId, &message.SentAt, &message.ReceivedAt)

		if err != nil {
			return nil, newErrRepository(err.Error())
		}

		geom, err := geos.FromHex(string(loc))

		if err != nil {
			return nil, newErrRepository(err.Error())
		}

		message.Location.Long, err = geom.X()

		if err != nil {
			return nil, newErrRepository(err.Error())
		}

		message.Location.Lat, err = geom.Y()

		if err != nil {
			return nil, newErrRepository(err.Error())
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func MakePostgresqlRespository(db *sql.DB) (MessageRepository, error) {
	return &postgresqlMessageRepository{db}, nil
}
