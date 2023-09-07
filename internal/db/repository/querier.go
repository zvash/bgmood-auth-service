// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package repository

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AttachRoleToUser(ctx context.Context, arg AttachRoleToUserParams) error
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error)
	DeleteExpiredTokens(ctx context.Context) error
	GetRoleByName(ctx context.Context, name string) (Role, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetSessionByAccessToken(ctx context.Context, accessToken string) (Session, error)
	GetSessionWithActiveRefreshToken(ctx context.Context, refreshToken string) (Session, error)
	GetTokenByToken(ctx context.Context, token string) (Token, error)
	GetTokenByTokenAndType(ctx context.Context, arg GetTokenByTokenAndTypeParams) (Token, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByAccessToken(ctx context.Context, accessToken string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUsersInfoByUserIds(ctx context.Context, userids []uuid.UUID) ([]User, error)
	ListActiveSessions(ctx context.Context, arg ListActiveSessionsParams) ([]Session, error)
	RegisterUser(ctx context.Context, arg RegisterUserParams) (User, error)
	TerminateOtherSessions(ctx context.Context, arg TerminateOtherSessionsParams) error
	TerminateSingleSession(ctx context.Context, arg TerminateSingleSessionParams) error
	UnsetUserAvatar(ctx context.Context, id uuid.UUID) (User, error)
	UpdateAccessToken(ctx context.Context, arg UpdateAccessTokenParams) (Session, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	VerifyEmail(ctx context.Context, id uuid.UUID) (User, error)
	VerifyEmailByEmail(ctx context.Context, email string) (User, error)
}

var _ Querier = (*Queries)(nil)
