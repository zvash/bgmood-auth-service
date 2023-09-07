package gapi

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/zvash/bgmood-auth-service/internal/authpb"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetUsersInfo(ctx context.Context, req *authpb.GetUsersInfoRequest) (*authpb.GetUsersInfoResponse, error) {
	dto := pbGetUsersInfoRequestToValGetUsersInfoRequest(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	_, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	userUUIDs := make([]uuid.UUID, 0)
	for _, userId := range dto.UserIds {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			continue
		}
		userUUIDs = append(userUUIDs, userUUID)
	}
	if len(userUUIDs) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "non of provided user ids are valid")
	}
	dbUsers, err := server.db.GetUsersInfoByUserIds(ctx, userUUIDs)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "no record were found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error.")
	}
	users := make([]*authpb.User, 0)
	for _, dbUser := range dbUsers {
		user := repoUserToProtobufUser(dbUser)
		users = append(users, user)
	}
	resp := &authpb.GetUsersInfoResponse{
		Users: users,
	}
	return resp, nil
}
