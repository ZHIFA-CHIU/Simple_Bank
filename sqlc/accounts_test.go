package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zhifaq/simple_bank/utils"
)

func createTestAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
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

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	account := createTestAccount(t)

	retrievedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedAccount)
	require.Equal(t, account.ID, retrievedAccount.ID)
	require.Equal(t, account.Owner, retrievedAccount.Owner)
	require.Equal(t, account.Balance, retrievedAccount.Balance)
	require.Equal(t, account.Currency, retrievedAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, retrievedAccount.CreatedAt, 0)
}

func TestUpdateAccount(t *testing.T) {
	account := createTestAccount(t)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: utils.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, account.Currency, updatedAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, updatedAccount.CreatedAt, 0)
}

func TestDeleteAccount(t *testing.T) {
	account := createTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedAccount)
}

func TestListAccounts(t *testing.T) {
	for range 10 {
		createTestAccount(t)
	}

	arg := GetAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
