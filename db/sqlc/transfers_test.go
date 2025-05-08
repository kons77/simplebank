package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {

	transfer, err := testStore.CreateTransfer(context.Background(), argTestTransferParams)
	require.NoError(t, err) // check that the error must be nil and will automatically fail the test if it’s not
	require.NotEmpty(t, transfer)

	require.Equal(t, argTestTransferParams.FromAccountID, transfer.FromAccountID)
	require.Equal(t, argTestTransferParams.ToAccountID, transfer.ToAccountID)
	require.Equal(t, argTestTransferParams.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

}

func TestDeleteTransfer(t *testing.T) {

	transfer, err := testStore.CreateTransfer(context.Background(), argTestTransferParams)
	require.NoError(t, err)
	err = testStore.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

}

func TestGetTransfer(t *testing.T) {
	transfer1, err := testStore.CreateTransfer(context.Background(), argTestTransferParams)
	require.NoError(t, err)

	transfer2, err := testStore.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer1)

	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
}
