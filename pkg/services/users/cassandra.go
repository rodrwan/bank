package users

import (
	"context"
	"fmt"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/gocql/gocql"
	pb "github.com/rodrwan/bank/pkg/pb/users"
)

// Query query params to get data from database.
type Query struct {
	ID       string
	Username string
	Email    string
}

// ReadUsersStore is the interface of data store of read operations for users.
type ReadUsersStore interface {
	Get(ctx context.Context, query Query) (*pb.User, error)
}

// NewReadUsersDatabase instantiate a new connection to cassandra cluster.
func NewReadUsersDatabase(cassandraDSN string) (ReadUsersStore, error) {
	cluster := gocql.NewCluster(cassandraDSN)
	cluster.Keyspace = "test"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &ReadUsersDatabase{
		cqlSession: session,
	}, nil
}

// ReadUsersDatabase wrapper to connect cassandra database for users.
type ReadUsersDatabase struct {
	cqlSession *gocql.Session

	sync.RWMutex
}

// Get gets a User.
func (rud *ReadUsersDatabase) Get(ctx context.Context, query Query) (*pb.User, error) {
	rud.RLock()
	defer rud.RUnlock()

	fmt.Println("Getting all Employees")

	q := squirrel.Select("id,first_name,last_name,email,username").From("users")

	if query.ID != "" {
		q = q.Where("id = ?", query.ID)
	}

	if query.Username != "" {
		q = q.Where("username = ?", query.Username)
	}

	if query.Email != "" {
		q = q.Where("email = ?", query.Email)
	}

	sql, args, err := q.PlaceholderFormat(squirrel.Question).ToSql()
	if err != nil {
		return nil, err
	}

	iter := rud.cqlSession.Query(sql, args...).Iter()
	var user *pb.User
	m := map[string]interface{}{}

	for iter.MapScan(m) {
		user = &pb.User{
			Id:        m["id"].(string),
			FirstName: m["first_name"].(string),
			LastName:  m["last_name"].(string),
			Email:     m["email"].(string),
			Username:  m["username"].(string),
		}
	}

	return user, nil
}
