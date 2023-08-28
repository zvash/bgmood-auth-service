package gapi

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/zvash/bgmood-auth-service/internal/authpb"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) TerminateSingleSession(
	ctx context.Context,
	req *pb.TerminateSingleSessionRequest,
) (*pb.TerminateSingleSessionResponse, error) {
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	sessionUUID, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "session id must be a valid uuid value")
	}
	session, err := server.db.GetSession(ctx, sessionUUID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "session was not found.")
	}
	if session.UserID != user.ID {
		return nil, status.Errorf(codes.NotFound, "session was not found.")
	}
	err = server.db.TerminateSingleSession(ctx, repository.TerminateSingleSessionParams{
		ID:     sessionUUID,
		UserID: user.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove the session")
	}
	resp := &pb.TerminateSingleSessionResponse{
		Message: "successfully removed the session",
	}
	return resp, nil
}
