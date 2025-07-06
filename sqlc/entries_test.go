package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zhifaq/simple_bank/utils"
)

func TestCreateEntry(t *testing.T) {
	account := createTestAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func createRandomEntry(t *testing.T) Entry {
	account := createTestAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	retrievedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedEntry)
	require.Equal(t, entry.ID, retrievedEntry.ID)
	require.Equal(t, entry.AccountID, retrievedEntry.AccountID)
	require.Equal(t, entry.Amount, retrievedEntry.Amount)
	require.WithinDuration(t, entry.CreatedAt, retrievedEntry.CreatedAt, 0)
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID:     entry.ID,
		Amount: utils.RandomMoney(),
	}

	updatedEntry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)
	require.Equal(t, entry.ID, updatedEntry.ID)
	require.Equal(t, entry.AccountID, updatedEntry.AccountID)
	require.Equal(t, arg.Amount, updatedEntry.Amount)
	require.WithinDuration(t, entry.CreatedAt, updatedEntry.CreatedAt, 0)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	retrievedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, retrievedEntry)
}

func TestGetEntries(t *testing.T) {
	for range 10 {
		createRandomEntry(t)
	}

	arg := GetEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.GetEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}