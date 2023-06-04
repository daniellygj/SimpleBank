package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	// This is composition. is to extend struct functionality instead of inheritance.
	// All individual query functions provided by Queries will be available to Store.
	*Queries

	//Required to create a new db transaction
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
// func(queries *Queries) Its a callback function
func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTX performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts´balance within a single database transaction
func (store *Store) TransferTX(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = createNewTransfer(q, ctx, params)
		if err != nil {
			return err
		}

		result.FromEntry, err = createNewEntry(q, ctx, params.FromAccountId, -params.Amount)
		if err != nil {
			return err
		}

		result.ToEntry, err = createNewEntry(q, ctx, params.ToAccountId, params.Amount)
		if err != nil {
			return err
		}

		result.FromAccount, err = updateAccountsBalance(q, ctx, params.FromAccountId, -params.Amount)
		if err != nil {
			return err
		}

		result.ToAccount, err = updateAccountsBalance(q, ctx, params.ToAccountId, params.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func updateAccountsBalance(q *Queries, ctx context.Context, accountId int64, amount int64) (Account, error) {
	return q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountId,
		Amount: amount,
	})
}

func createNewEntry(q *Queries, ctx context.Context, accountId int64, amount int64) (Entry, error) {
	return q.CreateEntry(ctx, CreateEntryParams{
		AccountID: accountId,
		Amount:    amount,
	})
}

func createNewTransfer(q *Queries, ctx context.Context, params TransferTxParams) (Transfer, error) {
	return q.CreateTransfer(ctx, CreateTransferParams{
		FromAccountID: params.FromAccountId,
		ToAccountID:   params.ToAccountId,
		Amount:        params.Amount,
	})
}
