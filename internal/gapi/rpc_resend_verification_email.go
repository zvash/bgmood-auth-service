package gapi

import (
	"context"
	"fmt"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ResendVerificationEmail(
	ctx context.Context,
	req *pb.ResendVerificationEmailRequest,
) (*pb.ResendVerificationEmailResponse, error) {
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	err = server.sendVerifyEmail(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	resp := &pb.ResendVerificationEmailResponse{
		Message: fmt.Sprintf("verification email was sent to: %s", user.Email),
	}
	return resp, nil
}
