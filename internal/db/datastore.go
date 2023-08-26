package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
)

type DataStore interface {
	repository.Querier
	RegisterUserTransaction(ctx context.Context, arg RegisterUserTransactionParams) (RegisterUserTransactionResult, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*repository.Queries
}

func NewDataStore(connPool *pgxpool.Pool) DataStore {
	return &SQLStore{
		connPool: connPool,
		Queries:  repository.New(connPool),
	}
}

func (store *SQLStore) executeTransaction(ctx context.Context, fn func(queries *repository.Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	q := repository.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
