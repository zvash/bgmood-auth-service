package db

import (
	"context"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
)

type RegisterUserTransactionParams struct {
	repository.RegisterUserParams
	AfterRegister func(user repository.User) error
}

type RegisterUserTransactionResult struct {
	User repository.User
}

func (store *SQLStore) RegisterUserTransaction(ctx context.Context, arg RegisterUserTransactionParams) (RegisterUserTransactionResult, error) {
	var result RegisterUserTransactionResult
	err := store.executeTransaction(ctx, func(queries *repository.Queries) error {
		var err error
		result.User, err = queries.RegisterUser(ctx, repository.RegisterUserParams{
			ID:       arg.ID,
			Name:     arg.Name,
			Email:    arg.Email,
			Password: arg.Password,
		})
		if err != nil {
			return err
		}

		if err := arg.AfterRegister(result.User); err != nil {
			return err
		}
		normalRole, err := queries.GetRoleByName(ctx, "normal")
		if err != nil {
			return err
		}
		err = queries.GiveRoleToUser(ctx, repository.GiveRoleToUserParams{
			UserID: result.User.ID,
			RoleID: normalRole.ID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
