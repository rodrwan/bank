package users

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/rodrwan/bank/pkg/pb/users"
	queuemanager "github.com/rodrwan/bank/pkg/queueManager"
)

// WriteServiceConfig write services configuration.
type WriteServiceConfig struct {
	Queue queuemanager.QueueManager
	Store WriteUsersStore
}

// NewWriteServer initialises an instance of the accounts write server.
func NewWriteServer(asc WriteServiceConfig) pb.UsersWriteServiceServer {
	fmt.Println("Accounts write service")

	return WriteService{
		store: asc.Store,
		queue: asc.Queue,
	}
}

// WriteService wrap all write operation.
type WriteService struct {
	store WriteUsersStore
	queue queuemanager.QueueManager

	pb.UnimplementedUsersWriteServiceServer
}

// CreateUser creates a new user.
func (svc WriteService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	fmt.Println("CreateUser...")
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	if err := svc.store.Create(ctx, in.User); err != nil {
		return nil, err
	}

	event := queuemanager.NewEvent(CreatedUserEvent, "user-svc", data)
	if err := svc.queue.Publish(event); err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{}, nil
}
