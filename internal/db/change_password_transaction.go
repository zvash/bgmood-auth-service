package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
)

func (store *SQLStore) ChangePasswordTransaction(
	ctx context.Context,
	user repository.User,
	hashedPassword string,
	shouldTerminate bool,
	currentSessionID uuid.UUID,
) (bool, error) {
	err := store.executeTransaction(ctx, func(queries *repository.Queries) error {
		_, err := queries.UpdateUser(ctx, repository.UpdateUserParams{
			ID:       user.ID,
			Password: pgtype.Text{String: hashedPassword, Valid: true},
		})
		if err != nil {
			return err
		}
		if shouldTerminate {
			err = queries.TerminateOtherSessions(ctx, repository.TerminateOtherSessionsParams{
				ID:     currentSessionID,
				UserID: user.ID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
