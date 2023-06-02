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

	for i := 0; i < transactionsQty; i++ {
		// Receive messages from the channel
		err := <-errors
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

		//	TODO: check accounts balance

	}

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
