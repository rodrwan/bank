package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/rodrwan/bank/pkg/auth"
	"github.com/rodrwan/bank/pkg/pb/accounts"
	"github.com/rodrwan/bank/pkg/pb/users"
	generated1 "github.com/rodrwan/bank/pkg/services/graph/graph/generated"
	model1 "github.com/rodrwan/bank/pkg/services/graph/graph/model"
	"github.com/twinj/uuid"
)

func (r *mutationResolver) CreateUserAccount(ctx context.Context, input model1.NewUser) (*model1.CreateUserAccountResponse, error) {
	user := &users.User{
		Id:        uuid.NewV4().String(),
		FirstName: input.Name,
	}

	account := &accounts.Account{
		UserId: user.Id,
	}
	if err := r.CommandHandler.CreateUserAccount(ctx, user, account); err != nil {
		return nil, err
	}

	return &model1.CreateUserAccountResponse{
		Status: true,
		UserID: user.Id,
	}, nil
}

func (r *mutationResolver) CreateSession(ctx context.Context, input model1.NewSession) (*model1.Session, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateDeposit(ctx context.Context, input model1.NewDeposit) (*model1.Deposit, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Profile(ctx context.Context) (*model1.Profile, error) {
	userID := auth.GetUserIDFromContext(ctx)

	user, account, err := r.QueryService.GetUserAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model1.Profile{
		User: &model1.User{
			ID:   userID,
			Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		},
		Account: &model1.Account{
			UserID:  account.UserId,
			Name:    account.Name,
			Number:  account.Number,
			Balance: fmt.Sprintf("%d", account.Balance),
		},
	}, nil
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated1.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
