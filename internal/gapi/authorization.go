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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, user, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, user, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, user, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, user, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	user, err := server.db.GetUserByAccessToken(ctx, accessToken)
	if err != nil {
		return nil, user, fmt.Errorf("invalid access token: %s", err)
	}
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, user, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, user, nil
}

func (server *Server) getAuthenticatedUser(ctx context.Context) (repository.User, error) {
	user := repository.User{}
	_, user, err := server.authorizeUser(ctx)
	if err != nil {
		return user, err
	}
	return user, nil
}
