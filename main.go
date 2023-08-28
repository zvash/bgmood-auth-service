package main

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/zvash/bgmood-auth-service/internal/authpb"
	"github.com/zvash/bgmood-auth-service/internal/db"
	"github.com/zvash/bgmood-auth-service/internal/gapi"
	"github.com/zvash/bgmood-auth-service/internal/util"
	"github.com/zvash/bgmood-auth-service/internal/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dataStore := createDBConnectionPool(config.DBSource)

	redisOpt := createRedisClientOption(config)
	messagePublisher := worker.NewRedisMessagePublisher(redisOpt)

	runGrpcServer(config, dataStore, messagePublisher)
}

func createDBConnectionPool(dbSource string) db.DataStore {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	res, err := connPool.Query(context.Background(), "SHOW TRANSACTION ISOLATION LEVEL;")
	if err != nil {
		log.Fatal("error connecting to db:", err)
	}
	res.Next()
	values, err := res.Values()
	if err != nil {
		return nil
	}
	fmt.Println(values)
	return db.NewDataStore(connPool)
}

func runGrpcServer(config util.Config, dataStore db.DataStore, messagePublisher worker.MessagePublisher) {
	server, err := gapi.NewServer(config, dataStore, messagePublisher)
	if err != nil {
		log.Fatal("cannot create gRPC server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func createRedisClientOption(config util.Config) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
}
