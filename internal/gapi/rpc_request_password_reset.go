package gapi

import (
	"context"
	"errors"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (server *Server) RequestPasswordReset(
	ctx context.Context,
	req *pb.RequestPasswordResetRequest,
) (*pb.RequestPasswordResetResponse, error) {
	dto := pbRequestPasswordResetRequestToValRequestPasswordResetRequest(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	_ = server.db.DeleteExpiredTokens(ctx)
	user, err := server.db.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return &pb.RequestPasswordResetResponse{
				Message: "password reset token was sent to the given email.", //🤥
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	token, err := server.createUniqueToken(ctx, user.Email, repository.TokenTypePASSWORDRESET, server.config.PasswordResetDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	err = server.sendPasswordResetEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	log.Println(token.Token)
	resp := &pb.RequestPasswordResetResponse{
		Message: "password reset token was sent to the given email.",
	}
	return resp, nil
}
