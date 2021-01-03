package users

import (
	"context"
	"fmt"

	pb "github.com/rodrwan/bank/pkg/pb/users"
)

// ReadServiceConfig read services configuration.
type ReadServiceConfig struct {
	Store ReadUsersStore
}

// NewReadService initialises an instance of the users service.
func NewReadService(asc ReadServiceConfig) pb.UsersReadServiceServer {
	fmt.Println("Users read service")
	return ReadService{
		store: asc.Store,
	}
}

// ReadService wrap all read operation.
type ReadService struct {
	store ReadUsersStore

	pb.UnimplementedUsersReadServiceServer
}

// GetUser gets a user by given queries.
func (svc ReadService) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := svc.store.Get(ctx, Query{
		ID:       in.Id,
		Email:    in.Email,
		Username: in.Username,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user,
	}, nil
}
