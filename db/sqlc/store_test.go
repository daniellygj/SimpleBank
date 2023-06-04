package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTX(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	//	run in concurrent transfer transactions
	transactionsQty := 5
	amount := int64(10)

	// channel is designed to connect concurrent Go routines
	// And allow them to safely share data with each other without explicit locking
	errors := make(chan error)
	results := make(chan TransferTxResult)
	existed := make(map[int]bool)

	for i := 0; i < transactionsQty; i++ {
		go func() {
			result, err := store.TransferTX(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			// send data to the channels
			errors <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < transactionsQty; i++ {
		err := <-errors // Receive messages from the channel
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		checkTransfer(t, result, account1, account2, amount, store)

		fromEntry := result.FromEntry
		fromEntryId := account1.ID
		err = checkEntry(t, fromEntry, fromEntryId, -amount, store)

		toEntry := result.ToEntry
		toEntryId := account2.ID
		err = checkEntry(t, toEntry, toEntryId, amount, store)

		checkAccounts(t, result, account1, account2, amount, transactionsQty, existed)

	}

	checkUpdatedBalance(t, account1, transactionsQty, amount, account2)
}

func checkUpdatedBalance(t *testing.T, account1 Account, transactionsQty int, amount int64, account2 Account) {
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-int64(transactionsQty)*amount, updatedAccount1.Balance)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance+int64(transactionsQty)*amount, updatedAccount2.Balance)

}

func checkAccounts(t *testing.T, result TransferTxResult, account1 Account, account2 Account, amount int64, transactionsQty int, existed map[int]bool) {
	// Check accounts
	fromAccount := result.FromAccount
	require.NotEmpty(t, fromAccount)
	require.Equal(t, fromAccount.ID, account1.ID)

	toAccount := result.ToAccount
	require.NotEmpty(t, toAccount)
	require.Equal(t, toAccount.ID, account2.ID)

	// check accounts balance
	diffFromAccount := account1.Balance - fromAccount.Balance
	diffToAccount := toAccount.Balance - account2.Balance
	require.Equal(t, diffToAccount, diffFromAccount)
	require.True(t, diffFromAccount > 0)
	require.True(t, diffFromAccount%amount == 0)

	k := int(diffFromAccount / amount)
	require.True(t, k >= 1 && k <= transactionsQty)
	require.NotContains(t, existed, k)
	existed[k] = true
}

func checkEntry(t *testing.T, entry Entry, entryId int64, amount int64, store *Store) error {
	require.NotEmpty(t, entry)
	require.Equal(t, entryId, entry.AccountID)
	require.Equal(t, amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	_, err := store.GetEntry(context.Background(), entryId)
	require.NoError(t, err)
	return err
}

func checkTransfer(t *testing.T, result TransferTxResult, account1 Account, account2 Account, amount int64, store *Store) {
	transfer := result.Transfer
	require.NotEmpty(t, transfer)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	_, err := store.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
}
