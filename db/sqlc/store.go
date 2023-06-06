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
func (store *Store) TransferTX(
	ctx context.Context,
	params TransferTxParams,
) (TransferTxResult, error) {
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

		err2 := updateToAndFromAccountsBalance(q, params, &result, ctx)
		if err2 != nil {
			return err2
		}

		return nil
	})

	return result, err
}

func updateToAndFromAccountsBalance(
	q *Queries,
	params TransferTxParams,
	result *TransferTxResult,
	ctx context.Context,
) (err error) {
	if params.FromAccountId < params.ToAccountId {
		result.FromAccount, result.ToAccount, err = addAmountToAccounts(ctx, q, params.FromAccountId, -params.Amount, params.ToAccountId, params.Amount)
	} else {
		result.ToAccount, result.FromAccount, err = addAmountToAccounts(ctx, q, params.ToAccountId, params.Amount, params.FromAccountId, -params.Amount)
	}
	return err
}

func addAmountToAccounts(
	ctx context.Context,
	q *Queries,
	accountId1 int64,
	amount1 int64,
	accountId2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = updateAccountBalance(q, ctx, accountId1, amount1)
	if err != nil {
		return
	}
	account2, err = updateAccountBalance(q, ctx, accountId2, amount2)
	return
}

func updateAccountBalance(
	q *Queries,
	ctx context.Context,
	accountId int64,
	amount int64,
) (Account, error) {
	return q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountId,
		Amount: amount,
	})
}

func createNewEntry(
	q *Queries,
	ctx context.Context,
	accountId int64,
	amount int64,
) (Entry, error) {
	return q.CreateEntry(ctx, CreateEntryParams{
		AccountID: accountId,
		Amount:    amount,
	})
}

func createNewTransfer(
	q *Queries,
	ctx context.Context,
	params TransferTxParams,
) (Transfer, error) {
	return q.CreateTransfer(ctx, CreateTransferParams{
		FromAccountID: params.FromAccountId,
		ToAccountID:   params.ToAccountId,
		Amount:        params.Amount,
	})
}
