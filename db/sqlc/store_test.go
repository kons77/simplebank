package db

import (
	"context"
	"github/kons77/simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranferTx(t *testing.T) {
	//store := NewStore()

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			/* Now we cannot just use testify require to check them right here because this function
			is running inside a different go routine from the one that our TestTransferTx function is running on,
			So there’s no guarantee that it will stop the whole test if a condition is not satisfied.
			The correct way to verify the error and result is to send them back
			to the main go routine that our test is running on,
			and check them from there.
			*/
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: util.ToPgInt8(account1.ID),
				ToAccountID:   util.ToPgInt8(account2.ID),
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check resultes
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check tranfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
		require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTrasfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID.Int64)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID.Int64)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: check accounts' balance
	}
}
