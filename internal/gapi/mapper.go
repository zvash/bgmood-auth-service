package gapi

import (
	"fmt"
	"github.com/iancoleman/strcase"
	pb "github.com/zvash/bgmood-auth-service/internal/authpb"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func errorResponseToErrorDetailsBadRequestFieldViolation(er val.ErrorResponse) *errdetails.BadRequest_FieldViolation {
	fieldName := strcase.ToSnake(er.FailedField)
	return &errdetails.BadRequest_FieldViolation{
		Field: fieldName,
		Description: fmt.Sprintf(
			"[%s]: '%v' | Needs to implement '%s'",
			fieldName,
			er.Value,
			er.Tag,
		),
	}
}

func errorResponsesToErrorDetailsBadRequestFieldViolations(ers []val.ErrorResponse) (violations []*errdetails.BadRequest_FieldViolation) {
	for _, er := range ers {
		violations = append(violations, errorResponseToErrorDetailsBadRequestFieldViolation(er))
	}
	return
}

func errorResponsesToStatusErrors(errs []val.ErrorResponse) error {
	violations := errorResponsesToErrorDetailsBadRequestFieldViolations(errs)
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")
	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}
	return status.Errorf(statusDetails.Code(), "%v", statusDetails.Details())
}

func pbRegisterUserRequestToValRegisterUserRequest(pbr *pb.RegisterUserRequest) val.RegisterUserRequest {
	return val.RegisterUserRequest{
		Email:                pbr.Email,
		Name:                 pbr.Name,
		Password:             pbr.Password,
		PasswordConfirmation: pbr.PasswordConfirmation,
	}
}

func repoUserToProtobufUser(repoUser repository.User) *pb.User {
	user := &pb.User{
		Id:         repoUser.ID.String(),
		Email:      repoUser.Email,
		Name:       repoUser.Name,
		Avatar:     nil,
		IsVerified: !repoUser.VerifiedAt.Time.IsZero(),
		CreatedAt:  timestamppb.New(repoUser.CreatedAt),
	}
	if repoUser.Avatar.Valid {
		user.Avatar = &repoUser.Avatar.String
	}
	return user
}

func pbLoginRequestToValLoginRequest(pbl *pb.LoginRequest) val.LoginRequest {
	return val.LoginRequest{
		Email:    pbl.GetEmail(),
		Password: pbl.GetPassword(),
	}
}

func pbChangePasswordRequestToValChangePasswordRequest(pbc *pb.ChangePasswordRequest) val.ChangePasswordRequest {
	return val.ChangePasswordRequest{
		CurrentPassword:         pbc.CurrentPassword,
		NewPassword:             pbc.NewPassword,
		NewPasswordConfirmation: pbc.NewPasswordConfirmation,
		TerminateOtherSessions:  pbc.TerminateOtherSessions,
	}
}

func pbRequestPasswordResetRequestToValRequestPasswordResetRequest(
	pbr *pb.RequestPasswordResetRequest,
) val.RequestPasswordResetRequest {
	return val.RequestPasswordResetRequest{
		Email: pbr.GetEmail(),
	}
}

func pbResetPasswordRequestToValResetPasswordRequest(pbr *pb.ResetPasswordRequest) val.ResetPasswordRequest {
	return val.ResetPasswordRequest{
		Email:                pbr.Email,
		Token:                pbr.Token,
		Password:             pbr.Password,
		PasswordConfirmation: pbr.PasswordConfirmation,
	}
}

func pbUpdateUserRequestToValUpdateUserRequest(pbu *pb.UpdateUserRequest) val.UpdateUserRequest {
	return val.UpdateUserRequest{
		Name: pbu.GetName(),
	}
}

func pbGetUsersInfoRequestToValGetUsersInfoRequest(req *pb.GetUsersInfoRequest) val.GetUsersInfoRequest {
	return val.GetUsersInfoRequest{
		UserIds: req.GetUserIds(),
	}
}
