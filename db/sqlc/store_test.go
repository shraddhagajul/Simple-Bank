package db

import (
	"context"
	"fmt"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	fmt.Println(">>before: ",fromAccount.Balance,toAccount.Balance)

	goRountineCount := 5
	transferAmount := int64(10)
	
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i:= 0;i < goRountineCount; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID: toAccount.ID,
				Amount: transferAmount,
			})
			errs <- err
			results <- result
		}()
	}
	
	existed := make(map[int]bool)
	//check results
	for i:= 0 ;i<goRountineCount;i++ {
			err := <-errs
			require.NoError(t,err)

			result := <-results
			require.NotEmpty(t,result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, transferAmount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -transferAmount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, transferAmount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		senderAccount := result.FromAccount
		require.NotEmpty(t,senderAccount)
		require.Equal(t,fromAccount.ID, senderAccount.ID)

		receiverAccount := result.ToAccount
		require.NotEmpty(t,receiverAccount)
		require.Equal(t,toAccount.ID,receiverAccount.ID)

		fmt.Println(">>tx: ",senderAccount.Balance,receiverAccount.Balance)
		//check accounts' balance
		diffSenderAccount := fromAccount.Balance - senderAccount.Balance
		diffReceiverAccount := receiverAccount.Balance - toAccount.Balance
		require.Equal(t,diffSenderAccount,diffReceiverAccount)
		require.True(t,diffSenderAccount>0)
		require.True(t,diffSenderAccount%transferAmount == 0)

		k := int(diffSenderAccount/transferAmount)
		fmt.Println("K : ",k)
		require.True(t, k>=1 && k<=goRountineCount)
		require.NotContains(t,existed,k)
		existed[k] = true
	}

	//check final updated balance
	updatedSenderAccount, err := testQueries.GetAccount(context.Background(),fromAccount.ID)
	require.NoError(t,err)

	updatedReceiverAccount, err := testQueries.GetAccount(context.Background(),toAccount.ID)
	require.NoError(t,err)

	fmt.Println(">>after: ",updatedSenderAccount.Balance,updatedReceiverAccount.Balance)

	require.Equal(t, fromAccount.Balance - int64(goRountineCount)*transferAmount, updatedSenderAccount.Balance)
	require.Equal(t, toAccount.Balance + int64(goRountineCount)*transferAmount, updatedReceiverAccount.Balance)
}