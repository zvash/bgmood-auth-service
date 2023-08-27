package gapi

import (
	"context"
	"errors"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

func (server *Server) getAllActiveSessions(ctx context.Context) ([]repository.Session, string, error) {
	accessToken, err := server.getAccessToken(ctx)
	if err != nil {
		return nil, "", status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, accessToken, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	sessions, err := server.db.ListActiveSessions(ctx, repository.ListActiveSessionsParams{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, accessToken, status.Errorf(codes.NotFound, "no sessions were found.")
		}
		return nil, accessToken, status.Errorf(codes.Internal, "internal server error.")
	}
	return sessions, accessToken, nil
}

func (server *Server) createUniqueToken(ctx context.Context, email string, duration time.Duration) (repository.Token, error) {
	for {
		token := strings.ToUpper(util.RandomAlphaNumString(6))
		record, err := server.db.CreateToken(ctx, repository.CreateTokenParams{
			Email:     email,
			Token:     token,
			ExpiresAt: time.Now().Add(duration),
		})
		if err != nil {
			if db.ErrorCode(err) == db.UniqueViolation {
				continue
			}
			return repository.Token{}, err
		}
		return record, nil
	}
}
