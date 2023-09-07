package gapi

import (
	"context"
	"errors"
	"github.com/zvash/bgmood-auth-service/internal/authpb"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	_, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	accessToken, err := server.getAccessToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	session, err := server.db.GetSessionByAccessToken(ctx, accessToken)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "session was not found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	err = server.db.TerminateSingleSession(ctx, repository.TerminateSingleSessionParams{
		ID:     session.ID,
		UserID: session.UserID,
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "session was not found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	return &authpb.LogoutResponse{}, nil
}
