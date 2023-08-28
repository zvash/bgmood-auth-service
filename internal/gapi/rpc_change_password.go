package gapi

import (
	"context"
	pb "github.com/zvash/bgmood-auth-service/internal/authpb"
	"github.com/zvash/bgmood-auth-service/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	dto := pbChangePasswordRequestToValChangePasswordRequest(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	err = util.CheckPassword(req.GetCurrentPassword(), user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	hashedPassword, err := util.HashPassword(dto.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}
	session, err := server.getSession(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to change user's password. Session was not found.")
	}

	_, err = server.db.ChangePasswordTransaction(ctx, user, hashedPassword, dto.TerminateOtherSessions, session.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to change user's password")
	}
	resp := &pb.ChangePasswordResponse{
		Message: "Successfully changed the password",
	}
	return resp, nil
}
