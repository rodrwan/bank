package accounts

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"
	pb "github.com/rodrwan/bank/pkg/pb/accounts"
)

// WriteAccountsStore is the interface of data store of write operation for account.
type WriteAccountsStore interface {
	Create(ctx context.Context, acc *pb.Account) error
}

// WriteAccountsDatabase wrapper to connect postgres database for accounts.
type WriteAccountsDatabase struct {
	db *sqlx.DB
	sync.RWMutex
}

// Create creates an account
func (wad *WriteAccountsDatabase) Create(ctx context.Context, acc *pb.Account) error {
	wad.RLock()
	defer wad.RUnlock()

	return nil
}
