package gapi

import (
	"context"
	"errors"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) ListActiveSessions(ctx context.Context, req *pb.ListActiveSessionsRequest) (*pb.ListActiveSessionsResponse, error) {
	accessToken, err := server.getAccessToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	sessions, err := server.db.ListActiveSessions(ctx, repository.ListActiveSessionsParams{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "no sessions were found.")
		}
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	resp := &pb.ListActiveSessionsResponse{}
	for _, session := range sessions {
		item := &pb.ActiveSession{
			ClientIp:  session.ClientIp,
			UserAgent: session.UserAgent,
			ExpiresAt: timestamppb.New(session.ExpiresAt),
		}
		if session.AccessToken == accessToken {
			(*resp).CurrentSession = item
		} else {
			(*resp).OtherSessions = append((*resp).OtherSessions, item)
		}
	}
	return resp, nil
}
