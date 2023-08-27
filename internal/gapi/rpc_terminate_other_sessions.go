package gapi

import (
	"context"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) TerminateOtherSessions(
	ctx context.Context,
	req *pb.TerminateOtherSessionsRequest,
) (*pb.TerminateOtherSessionsResponse, error) {
	session, err := server.getSession(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	err = server.db.TerminateOtherSessions(ctx, repository.TerminateOtherSessionsParams{
		ID:     session.ID,
		UserID: session.UserID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to terminate other sessions")
	}
	resp := &pb.TerminateOtherSessionsResponse{
		Message: "successfully terminated all other sessions",
	}
	return resp, nil
}
