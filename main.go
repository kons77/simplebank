package main

import (
	"context"
	"log"

	"github.com/kons77/simplebank/api"
	db "github.com/kons77/simplebank/db/sqlc"
	"github.com/kons77/simplebank/util"

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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
