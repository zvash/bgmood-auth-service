package gapi

import (
	"context"
	"errors"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	payload, err := server.tokenMaker.VerifyToken(req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(payload.UserID, server.config.AccessTokenDuration)
	session, err := server.db.RefreshTokenTransaction(ctx, req.RefreshToken, accessToken)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
		}
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	resp := &pb.RefreshTokenResponse{
		AccessToken:          session.AccessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	}
	return resp, nil
}
