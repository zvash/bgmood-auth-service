package db

import (
	"context"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
)

func (store *SQLStore) RefreshTokenTransaction(ctx context.Context, refreshToken, newAccessToken string) (repository.Session, error) {
	var session repository.Session
	err := store.executeTransaction(ctx, func(queries *repository.Queries) error {
		currentSession, err := queries.GetSessionWithActiveRefreshToken(ctx, refreshToken)
		if err != nil {
			return err
		}
		session, err = queries.UpdateAccessToken(ctx, repository.UpdateAccessTokenParams{
			ID:          currentSession.ID,
			AccessToken: newAccessToken,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return session, err
}
