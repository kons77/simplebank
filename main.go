package main

import (
	"context"
	"github/kons77/simplebank/api"
	db "github/kons77/simplebank/db/sqlc"
	"github/kons77/simplebank/util"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var dbSource string
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		dbSource = config.DBSourceGH // path for GitHub Actions
	}
	dbSource = config.DBSourceLinux // path to the local linux machine

	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
