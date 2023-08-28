package gapi

import (
	"context"
	"errors"
	"github.com/hibiken/asynq"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/util"
	"github.com/zvash/bgmood-auth-service/internal/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

func (server *Server) getAllActiveSessions(ctx context.Context) ([]repository.Session, string, error) {
	accessToken, err := server.getAccessToken(ctx)
	if err != nil {
		return nil, "", status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, accessToken, status.Errorf(codes.Unauthenticated, "unauthorized.")
	}
	sessions, err := server.db.ListActiveSessions(ctx, repository.ListActiveSessionsParams{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, accessToken, status.Errorf(codes.NotFound, "no sessions were found.")
		}
		return nil, accessToken, status.Errorf(codes.Internal, "internal server error.")
	}
	return sessions, accessToken, nil
}

func (server *Server) createUniqueToken(ctx context.Context, email string, tokenType repository.TokenType, duration time.Duration) (repository.Token, error) {
	for {
		token := strings.ToUpper(util.RandomAlphaNumString(6))
		record, err := server.db.CreateToken(ctx, repository.CreateTokenParams{
			Email:     email,
			Token:     token,
			Type:      tokenType,
			ExpiresAt: time.Now().Add(duration),
		})
		if err != nil {
			if db.ErrorCode(err) == db.UniqueViolation {
				continue
			}
			return repository.Token{}, err
		}
		return record, nil
	}
}

func (server *Server) sendVerifyEmail(ctx context.Context, user repository.User) error {
	tokenRecord, err := server.createUniqueToken(ctx, user.Email, repository.TokenTypeVERIFICATIONEMAIL, server.config.VerifyEmailDuration)
	if err != nil {
		return err
	}
	payload := worker.PayloadSendVerifyEmail{
		Email: user.Email,
		Name:  user.Name,
		Token: tokenRecord.Token,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueDefault),
	}
	err = server.messagePublisher.PublishTaskSendVerifyEmail(ctx, &payload, opts...)
	return err
}

func (server *Server) sendPasswordResetEmail(ctx context.Context, email string) error {
	tokenRecord, err := server.createUniqueToken(ctx, email, repository.TokenTypePASSWORDRESET, server.config.PasswordResetDuration)
	if err != nil {
		return err
	}
	payload := worker.PayloadSendResetPasswordEmail{
		Email:   tokenRecord.Email,
		AppName: server.config.AppName,
		Token:   tokenRecord.Token,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueDefault),
	}
	err = server.messagePublisher.PublishTaskSendResetPasswordEmail(ctx, &payload, opts...)
	return err
}
