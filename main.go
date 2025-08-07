package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kons77/simplebank/api"
	db "github.com/kons77/simplebank/db/sqlc"
	"github.com/kons77/simplebank/gapi"
	"github.com/kons77/simplebank/pb"
	"github.com/kons77/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed doc/**
var embeddedFiles embed.FS

func main() {

	// load application configuration from the current directory
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := config.DBSource

	// create a new PostgreSQL connection pool using pgx driver
	// connection pools manage multiple database connections efficiently
	// and handle connection reuse, timeouts, and limits
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	// initialize the database store with the connection pool
	// this store likely contains methods for database operations (CRUD)
	store := db.NewStore(connPool)

	// start the gateway server in a separate goroutine (non-blocking)
	// this allows both servers to run concurrently
	go runGatewayServer(config, store)

	// start the gRPC server on the main goroutine (blocking)
	// this keeps the main function alive and handles gRPC requests
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

// start grpc gateway server
func runGatewayServer(config util.Config, store *db.SQLStore) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// to generate JSON with snake_case iunstead of camelCase
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// start fileserver with swagger documentation
	// --- fs := http.FileServer(http.Dir("./doc/swagger"))
	swaggerFS, err := fs.Sub(embeddedFiles, "doc/swagger")
	if err != nil {
		log.Fatal("failed to get embedded swagger folder:", err)
	}

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerFS))))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
	}
}

// start gin server -- legacy, for example
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
