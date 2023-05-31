package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simple_bank/util"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	account := createRandomAccount(t)
	accountFound, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountFound)

	require.Equal(t, account.ID, accountFound.ID)
	require.Equal(t, account.Owner, accountFound.Owner)
	require.Equal(t, account.Balance, accountFound.Balance)
	require.Equal(t, account.Currency, accountFound.Currency)
	require.WithinDuration(t, account.CreatedAt, accountFound.CreatedAt, time.Second)
}

func TestQueries_UpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	accountToUpdate := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	accountUpdated, err := testQueries.UpdateAccount(context.Background(), accountToUpdate)

	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, account.ID, accountUpdated.ID)
	require.Equal(t, account.Owner, accountUpdated.Owner)
	require.Equal(t, accountToUpdate.Balance, accountUpdated.Balance)
	require.Equal(t, account.Currency, accountUpdated.Currency)
	require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
}

func TestQueries_DeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	accountFound, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountFound)
}

func TestQueries_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accountsCreated := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accountListReturned, err := testQueries.ListAccounts(context.Background(), accountsCreated)
	require.NoError(t, err)
	require.Len(t, accountListReturned, 5)

	for _, account := range accountListReturned {
		require.NotEmpty(t, account)
	}
}
