package users

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"
	pb "github.com/rodrwan/bank/pkg/pb/users"
)

// WriteUsersStore is the interface of data store of write operations for users.
type WriteUsersStore interface {
	Create(ctx context.Context, user *pb.User) error
}

// NewWriteUsersDatabase instantiate a new postgres connection.
func NewWriteUsersDatabase(postgresDSN string) (WriteUsersStore, error) {
	db, err := sqlx.Open("postgres", postgresDSN)
	if err != nil {
		failOnError(err, "could not connect to db")
	}

	db.SetConnMaxLifetime(-1)

	return &WriteUsersDatabase{
		db: db,
	}, nil
}

// WriteUsersDatabase wrapper to connect postgres database for users.
type WriteUsersDatabase struct {
	db *sqlx.DB
	sync.RWMutex
}

// Create creates a new User.
func (wad *WriteUsersDatabase) Create(ctx context.Context, user *pb.User) error {
	wad.RLock()
	defer wad.RUnlock()

	return nil
}
