package gapi

import (
	"context"
	"fmt"
	"github.com/zvash/bgmood-auth-service/internal/db/repository"
	"github.com/zvash/bgmood-auth-service/internal/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, repository.User, error) {
	var user repository.User
	accessToken, err := server.getAccessToken(ctx)
	if err != nil {
		return nil, user, err
	}
	user, err = server.db.GetUserByAccessToken(ctx, accessToken)
	if err != nil {
		return nil, user, fmt.Errorf("invalid access token: %s", err)
	}
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, user, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, user, nil
}

func (server *Server) getAccessToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return "", fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return "", fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	return accessToken, nil
}

func (server *Server) getAuthenticatedUser(ctx context.Context) (repository.User, error) {
	user := repository.User{}
	_, user, err := server.authorizeUser(ctx)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (server *Server) getSession(ctx context.Context) (repository.Session, error) {
	var session repository.Session
	accessToken, err := server.getAccessToken(ctx)
	if err != nil {
		return session, err
	}
	return server.db.GetSessionByAccessToken(ctx, accessToken)
}
