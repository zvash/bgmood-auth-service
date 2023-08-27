package gapi

import (
	"context"
	"errors"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
