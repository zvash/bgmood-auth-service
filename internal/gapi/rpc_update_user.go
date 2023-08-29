package gapi

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zvash/bgmood-auth-service/internal/authpb"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *authpb.UpdateUserRequest) (*authpb.UpdateUserResponse, error) {
	dto := pbUpdateUserRequestToValUpdateUserRequest(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	user, err = server.db.UpdateUser(ctx, repository.UpdateUserParams{
		ID:     user.ID,
		Name:   pgtype.Text{String: req.GetName(), Valid: req.Name != nil},
		Avatar: pgtype.Text{String: req.GetAvatar(), Valid: req.Avatar != nil},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user's data")
	}
	resp := &authpb.UpdateUserResponse{
		User: repoUserToProtobufUser(user),
	}
	return resp, nil
}
