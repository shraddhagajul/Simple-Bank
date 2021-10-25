package db

import (
	"context"
	"database/sql"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
 }

 account, err := testQueries.CreateAccount(context.Background(),arg)
 require.NoError(t,err)
 require.NotEmpty(t,account)

 require.Equal(t, arg.Owner,account.Owner)
 require.Equal(t, arg.Balance,account.Balance)
 require.Equal(t, arg.Currency,account.Currency)

 require.NotZero(t,account.ID)
 require.NotZero(t,account.CreatedAt)

 return account
}
func TestCreateAccount(t *testing.T) {
		createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	accountInDb, err := testQueries.GetAccount(context.Background(),newAccount.ID)
	require.NoError(t,err)
	require.NotEmpty(t,accountInDb)

	require.Equal(t,newAccount.ID,accountInDb.ID)
	require.Equal(t,newAccount.Owner,accountInDb.Owner)
	require.Equal(t,newAccount.Balance,accountInDb.Balance)
	require.Equal(t,newAccount.Currency,accountInDb.Currency)

	require.WithinDuration(t,newAccount.CreatedAt,accountInDb.CreatedAt,time.Second)
	
}

func TestUpdateAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	arg := UpdateAccountParams{
		 ID: newAccount.ID,
		 Balance: util.RandomMoney(),
	}

	updatedRecord, err := testQueries.UpdateAccount(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,updatedRecord)

	require.Equal(t,newAccount.ID,updatedRecord.ID)
	require.Equal(t,newAccount.Owner,updatedRecord.Owner)
	require.Equal(t,arg.Balance,updatedRecord.Balance)
	require.Equal(t,newAccount.Currency,updatedRecord.Currency)
	require.WithinDuration(t,newAccount.CreatedAt,updatedRecord.CreatedAt,time.Second)

}

func TestDeleteAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(),newAccount.ID)
	require.NoError(t,err)

	getRecord, err := testQueries.GetAccount(context.Background(),newAccount.ID)

	require.Error(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,getRecord)
}

func TestListAccounts(t *testing.T) {
	for i:=0; i<10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	listAccounts, err := testQueries.ListAccounts(context.Background(),arg)
	require.NoError(t,err)
	require.Len(t,listAccounts,5)

	for _, account := range listAccounts {
		require.NotEmpty(t,account)
	}
}