package accounts

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/rodrwan/bank/pkg/pb/accounts"
)

// NewReadService initialises an instance of the accounts service.
func NewReadService(asc AccountServiceConfig) pb.AccountReadServiceServer {
	return ReadService{}
}

// ReadService ...
type ReadService struct {
	pb.UnimplementedAccountReadServiceServer
}

// GetAccountByUserID ...
func (svc ReadService) GetAccountByUserID(ctx context.Context, in *pb.GetAccountByUserIDRequest) (*pb.GetAccountByUserIDResponse, error) {
	inData, err := json.Marshal(in)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// get data from database
	var data *pb.Account
	if err := json.Unmarshal(inData, &data); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &pb.GetAccountByUserIDResponse{
		Account: data,
	}, nil
}
