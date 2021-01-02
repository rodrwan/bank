package graph

import (
	"context"

	accounts "github.com/rodrwan/bank/pkg/pb/accounts"
	users "github.com/rodrwan/bank/pkg/pb/users"
)

// ICommandHandler ...
type ICommandHandler interface {
	CreateUserAccount(ctx context.Context, user *users.User, account *accounts.Account) error
}

// CommandHandler ...
type CommandHandler struct {
	AccountsWriteClient accounts.AccountWriteServiceClient
	UsersWriteClient    users.UsersWriteServiceClient
}

// CreateUserAccount ...
func (ch *CommandHandler) CreateUserAccount(ctx context.Context, user *users.User, account *accounts.Account) error {
	_, err := ch.UsersWriteClient.CreateUser(ctx, &users.CreateUserRequest{User: user})
	if err != nil {
		return err
	}

	if _, err := ch.AccountsWriteClient.CreateAccount(ctx, &accounts.CreateAccountRequest{
		Account: account,
	}); err != nil {
		return err
	}

	return nil
}
