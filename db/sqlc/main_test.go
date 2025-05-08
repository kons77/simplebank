package db

import (
	"context"
	"github/kons77/simplebank/util"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var dbSource string
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		log.Println("Running inside GitHub Actions")
		dbSource = config.DBSourceGH // path for GitHub Actions
	} else {
		dbSource = config.DBSourceLinux // path to the local linux machine
	}

	// pgx.Connect - create only on connection
	testConnPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testStore = NewStore(testConnPool)

	os.Exit(m.Run())
}
