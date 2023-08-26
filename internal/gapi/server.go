package gapi

import (
	"fmt"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"github.com/zvash/bgmood-auth-service/internal/token"
	"github.com/zvash/bgmood-auth-service/internal/util"
	"github.com/zvash/bgmood-auth-service/internal/val"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedAuthServer
	config     util.Config
	db         db.DataStore
	validator  *val.XValidator
	tokenMaker token.Maker
}

func NewServer(config util.Config, dataStore db.DataStore) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		db:         dataStore,
		validator:  val.NewValidator(),
		tokenMaker: tokenMaker,
	}

	return server, nil
}
