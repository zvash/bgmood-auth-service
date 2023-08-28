package gapi

import (
	"context"
	"github.com/zvash/bgmood-auth-service/internal/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) Authenticate(ctx context.Context, req *authpb.AuthenticateRequest) (*authpb.AuthenticateResponse, error) {
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	resp := &authpb.AuthenticateResponse{
		User: repoUserToProtobufUser(user),
	}
	return resp, nil
}
