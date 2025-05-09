package main

import (
	"context"
	"log"
	"os"

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

	var dbSource string
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		dbSource = config.DBSourceGH // path for GitHub Actions
	} else {
		dbSource = config.DBSourceLinux // path to the local linux machine
	}

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
