package graph

import (
	"context"
	"errors"

	accounts "github.com/rodrwan/bank/pkg/pb/accounts"
	users "github.com/rodrwan/bank/pkg/pb/users"
)

// IQueryService ...
type IQueryService interface {
	GetUserAccount(ctx context.Context, userID string) (*users.User, *accounts.Account, error)
	GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (*users.User, error)
}

// QueryService ...
type QueryService struct {
	UsersReadClient    users.UsersReadServiceClient
	AccountsReadClient accounts.AccountReadServiceClient
}

// GetUserAccount ...
func (qs *QueryService) GetUserAccount(ctx context.Context, userID string) (*users.User, *accounts.Account, error) {
	userResp, err := qs.UsersReadClient.GetUser(ctx, &users.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, nil, err
	}

	accResp, err := qs.AccountsReadClient.GetAccountByUserID(ctx, &accounts.GetAccountByUserIDRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, nil, err
	}

	return userResp.User, accResp.Account, nil
}

// GetUserByUsernameAndPassword ...
func (qs *QueryService) GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (*users.User, error) {
	return nil, errors.New("not implemented")
}
