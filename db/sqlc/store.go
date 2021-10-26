package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store provides all functions to execute db queries and transactions
type Store struct {
	//WE can access all methods from store using queries for accessing them. Just like inheritance
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		//New() : creates and returns Queries object
		Queries: New(db),
	}
}

//execTx excutes a function within database transaction

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	//Start a new tx
	transaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//Create queries with that tx
	newQueriesObject := New(transaction)
	err = fn(newQueriesObject)
	if err != nil {
		//call the callback function with created queries
		if rbError := transaction.Rollback(); rbError != nil {
			return fmt.Errorf("tx err : %v, rbError : %v", err, rbError)
		}
		return err
	}

	//commit the changes
	return transaction.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id`
	ToAccountID   int64 `json:"to_account_id`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//TransferTx performs money transfer from one account to another
//It creates a transfer record, add account entries, update accounts' entries within a single transaction call
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// update accounts' balance
		fromAccount, err := q.GetAccountUpdated(ctx,arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx,UpdateAccountParams{
			  ID: fromAccount.ID,
				Balance: fromAccount.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}
		toAccount, err := q.GetAccountUpdated(ctx,arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx,UpdateAccountParams{
			ID: toAccount.ID,
			Balance: toAccount.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result,err
}
