package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zhifaq/simple_bank/utils"
)

func TestCreateTransfer(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	retrievedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedTransfer)
	require.Equal(t, transfer.ID, retrievedTransfer.ID)
	require.Equal(t, transfer.FromAccountID, retrievedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, retrievedTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, retrievedTransfer.Amount)
	require.WithinDuration(t, transfer.CreatedAt, retrievedTransfer.CreatedAt, 0)
}

func TestUpdateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     transfer.ID,
		Amount: utils.RandomMoney(),
	}

	updatedTransfer, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTransfer)
	require.Equal(t, transfer.ID, updatedTransfer.ID)
	require.Equal(t, transfer.FromAccountID, updatedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, updatedTransfer.ToAccountID)
	require.Equal(t, arg.Amount, updatedTransfer.Amount)
	require.WithinDuration(t, transfer.CreatedAt, updatedTransfer.CreatedAt, 0)
}

func TestGetTransfers(t *testing.T) {
	for range 10 {
		createRandomTransfer(t)
	}

	arg := GetTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.GetTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	retrievedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, retrievedTransfer)
}
