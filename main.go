package main

import (
	"context"
	"log"
	"net"

	"github.com/kons77/simplebank/api"
	db "github.com/kons77/simplebank/db/sqlc"
	"github.com/kons77/simplebank/gapi"
	"github.com/kons77/simplebank/pb"
	"github.com/kons77/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := config.DBSource

	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(connPool)
	runGrpcServer(config, store)

}

// start gRPC server
func runGrpcServer(config util.Config, store *db.SQLStore) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer) // allows the gRPC client to easily explore what RPCs are available on the server, and how to call them

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

// start gin server
func runGinServer(config util.Config, store *db.SQLStore) {

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
