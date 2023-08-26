package gapi

import (
	"context"
	"github.com/google/uuid"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"github.com/zvash/bgmood-auth-service/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (server *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	dto := pbRegisterUserRequestToValRegisterUserRequest(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	hashedPassword, err := util.HashPassword(dto.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not assign a new uuid: %v", err)
	}
	registerUserTransactionParams := db.RegisterUserTransactionParams{
		RegisterUserParams: repository.RegisterUserParams{
			ID:       randomUUID,
			Email:    dto.Email,
			Name:     dto.Name,
			Password: hashedPassword,
		},
		AfterRegister: func(user repository.User) error {
			//request notification service to send a verification email to the created user
			log.Printf("request notification service to send a verification email to the created user: email: %s\n", user.Email)
			return nil
		},
	}
	transactionResult, err := server.db.RegisterUserTransaction(ctx, registerUserTransactionParams)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "this user is already exists!")
		}
		return nil, status.Errorf(codes.Internal, "error while trying to create the new user.")
	}
	resp := &pb.RegisterUserResponse{
		User: repoUserToProtobufUser(transactionResult.User),
	}
	return resp, nil
}
