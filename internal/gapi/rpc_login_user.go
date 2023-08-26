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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	dto := pbLoginRequestToValLoginRequest(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	user, err := server.db.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Errorf(codes.Internal, "error while trying to log user in")
	}
	err = util.CheckPassword(req.GetPassword(), user.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID.String(), server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while trying to log user in")
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID.String(),
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while trying to log user in")
	}
	metadata := server.extractMetadata(ctx)
	userAgent := metadata.UserAgent
	ipAddress := metadata.ClientIP
	session, err := server.db.CreateSession(ctx, repository.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     ipAddress,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	loginResponse := &pb.LoginResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  repoUserToProtobufUser(user),
	}

	return loginResponse, nil
}
