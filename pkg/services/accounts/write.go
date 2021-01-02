package accounts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	pb "github.com/rodrwan/bank/pkg/pb/accounts"
	queuemanager "github.com/rodrwan/bank/pkg/queueManager"
)

// NewWriteServer initialises an instance of the accounts write server.
func NewWriteServer(asc AccountServiceConfig) pb.AccountWriteServiceServer {
	fmt.Println("Accounts write service")

	db, err := sqlx.Open("postgres", asc.PostgresDNS)
	if err != nil {
		failOnError(err, "could not connect to db")
	}

	db.SetConnMaxLifetime(-1)

	store := &WriteAccountsDatabase{
		db: db,
	}

	qq, err := queuemanager.MakeQueueManager(asc.RabbitMQURL)
	if err != nil {
		failOnError(err, "could not connect to queue")
	}

	return WriteService{
		store: store,
		queue: qq,
	}
}

// WriteService ...
type WriteService struct {
	store WriteAccountsStore
	queue queuemanager.QueueManager

	pb.UnimplementedAccountWriteServiceServer
}

// CreateAccount creates a new account
func (svc WriteService) CreateAccount(ctx context.Context, in *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	fmt.Println("CreateAccount...")
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	if err := svc.store.Create(ctx, in.Account); err != nil {
		return nil, err
	}

	event := queuemanager.NewEvent(CreatedAccountEvent, "account-svc", data)
	if err := svc.queue.Publish(event); err != nil {
		return nil, err
	}

	// here implements logit to create bank account
	return &pb.CreateAccountResponse{}, nil
}
