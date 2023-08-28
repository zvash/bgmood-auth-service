package gapi

import (
	"fmt"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/pb"
	"github.com/zvash/bgmood-auth-service/internal/token"
	"github.com/zvash/bgmood-auth-service/internal/util"
	"github.com/zvash/bgmood-auth-service/internal/val"
	"github.com/zvash/bgmood-auth-service/internal/worker"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedAuthServer
	config           util.Config
	db               db.DataStore
	validator        *val.XValidator
	tokenMaker       token.Maker
	messagePublisher worker.MessagePublisher
}

func NewServer(config util.Config, dataStore db.DataStore, taskDistributor worker.MessagePublisher) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:           config,
		db:               dataStore,
		validator:        val.NewValidator(),
		tokenMaker:       tokenMaker,
		messagePublisher: taskDistributor,
	}

	return server, nil
}
