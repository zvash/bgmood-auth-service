package gapi

import (
	"context"
	"errors"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"github.com/zvash/bgmood-auth-service/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	dto := pbResetPasswordRequestToValResetPasswordRequest(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	tokenRecord, err := server.db.GetTokenByTokenAndType(ctx, repository.GetTokenByTokenAndTypeParams{
		Token: req.Token,
		Type:  repository.TokenTypePASSWORDRESET,
	})
	if err != nil || tokenRecord.Email != req.Email {
		return nil, status.Errorf(codes.InvalidArgument, "verification code is not valid.")
	}

	user, err := server.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "couldn't find the user")
		}
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}
	_, err = server.db.ChangePasswordTransaction(ctx, user, hashedPassword, false, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reset user's password")
	}
	resp := &pb.ResetPasswordResponse{
		Message: "Successfully reset the password",
	}
	return resp, nil
}
