package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	// CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

// SQLStore provides all functions to execute db queries and trasactions
type SQLStore struct {
	// composition is to extend struct functionality in Golang instead of inheritance
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	ToAccountID   pgtype.Int8 `json:"to_account_id"`
	// FromAccountID int64 `json:"from_account_id"`
	//ToAccountID   int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccountID Account  `json:"from_account"`
	ToAccountID   Account  `json:"to_account"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' ballance within a single db transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTrasfer(ctx, CreateTrasferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // negative because money is moving out of this account
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // positive because money is moving in to this account
		})
		if err != nil {
			return err
		}

		// TODO: update accounts' ballance

		return nil
	})

	return result, err
}
