package gapi

import (
	"context"
	"fmt"
	pb "github.com/zvash/bgmood-auth-service/internal/authpb"
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
	if !user.VerifiedAt.Time.IsZero() {
		return nil, status.Errorf(codes.AlreadyExists, "your email address is already verified.")
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
