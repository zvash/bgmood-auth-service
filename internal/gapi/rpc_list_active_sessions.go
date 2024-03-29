package gapi

import (
	"context"
	pb "github.com/zvash/bgmood-auth-service/internal/authpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) ListActiveSessions(ctx context.Context, req *pb.ListActiveSessionsRequest) (*pb.ListActiveSessionsResponse, error) {
	sessions, accessToken, err := server.getAllActiveSessions(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListActiveSessionsResponse{}
	for _, session := range sessions {
		item := &pb.ActiveSession{
			SessionId: session.ID.String(),
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
