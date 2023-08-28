package gapi

import (
	"context"
	"errors"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	record, err := server.db.GetTokenByTokenAndType(ctx, repository.GetTokenByTokenAndTypeParams{
		Token: req.Token,
		Type:  repository.TokenTypeVERIFICATIONEMAIL,
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, "verification code is not valid.")
		}
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	if user.Email != record.Email {
		return nil, status.Errorf(codes.InvalidArgument, "verification code is not valid.")
	}
	_, err = server.db.VerifyEmail(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	resp := &pb.VerifyEmailResponse{
		Message: "email address successfully verified",
	}
	return resp, nil
}
