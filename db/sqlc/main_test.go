package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kons77/simplebank/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {

	os.Setenv("TEST_ENV", "true")

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := config.DBSource
	if dbSource == "" {
		log.Fatal("DB_SOURCE environment variable not set")
	}

	for _, key := range []string{"DB_SOURCE", "TOKEN_SYMMETRIC_KEY"} {
		fmt.Printf("%s = %s\n", key, os.Getenv(key))
	}

	// pgx.Connect - create only on connection
	testConnPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testStore = NewStore(testConnPool)

	os.Exit(m.Run())
}
