package service

import (
	"database/sql"
	"errors"
	"time"
)

// MessageRepository represents a data source through which users can be managed.
type MessageRepository interface {
	AddMessage(message Message) (StoredMessage, error)
	GetMessage(id string) (StoredMessage, error)
	GetMessagesForLocation(location Location, radiusMeters float64, limit int, after time.Time) ([]StoredMessage, error)
}

// NewMessageRepository constructs a UserRepository from the given configuration.
func NewMessageRepository(config Configuration) (MessageRepository, error) {
	var err error
	var repo MessageRepository
	//var db *sql.DB
	switch config.GetRepoType() {
	case InMemoryRepo:
		repo, err = MakeInMemoryRepository(config)
	case PostgreSqlRepo:
		db, err := sql.Open("postgres", config.GetPgUrl())

		if err != nil {
			return nil, err
		}
		repo, err = MakePostgresqlRespository(db)
	default:
		err = newErrRepository("repository type unimplemented")
	}

	return repo, err
}

type errRepository struct {
	err error
}

func (er errRepository) Error() string {
	return er.err.Error()
}

func newErrRepository(msg string) error {
	return errRepository{errors.New(msg)}
}
