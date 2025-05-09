package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/kons77/simplebank/util"

	"github.com/stretchr/testify/require"
)

func TestTranferTx(t *testing.T) {
	//store := NewStore()

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			/* Now we cannot just use testify require to check them right here because this function
			is running inside a different go routine from the one that our TestTransferTx function is running on,
			So there’s no guarantee that it will stop the whole test if a condition is not satisfied.
			The correct way to verify the error and result is to send them back
			to the main go routine that our test is running on,
			and check them from there.
			*/
			//ctx := context.WithValue(context.Background(), txKey, txName)
			ctx := context.Background()
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: util.ToPgInt8(account1.ID),
				ToAccountID:   util.ToPgInt8(account2.ID),
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check resultes
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check tranfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID.Valid)
		require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
		require.True(t, transfer.ToAccountID.Valid)
		require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.True(t, fromEntry.AccountID.Valid)
		require.Equal(t, account1.ID, fromEntry.AccountID.Int64)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.True(t, toEntry.AccountID.Valid)
		require.Equal(t, account2.ID, toEntry.AccountID.Int64)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts' balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0) // amount, 2*amount, 3*amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// fmt.Println(">> after:", updatedAccount1.Balance, account2.Balance)
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTranferTxDeadlock(t *testing.T) {

	/*The idea is to have 5 transactions that send money from account 1 to account 2,
	and another 5 transactions that send money in reverse direction from account 2 to account 1.
	In this case, we only need to check for deadlock error.
	We don’t need to care about the result because it has already been checked in the other test*/

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		FromAccountID := account1.ID
		ToAccountID := account2.ID

		if i%2 == 1 {
			FromAccountID = account2.ID
			ToAccountID = account1.ID
		}

		go func() {

			ctx := context.Background()
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: util.ToPgInt8(FromAccountID),
				ToAccountID:   util.ToPgInt8(ToAccountID),
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// check resultes
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// fmt.Println(">> after:", updatedAccount1.Balance, account2.Balance)
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
