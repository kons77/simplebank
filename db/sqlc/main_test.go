package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
)

var dbSource = getDBSource()

func getDBSource() string {
	if os.Getenv("CI") == "true" {
		// GitHub Actions
		return "postgresql://postgres:secret@localhost:5438/simplebank?sslmode=disable"
	}
	// Locally
	return "postgresql://postgres:secret@192.168.88.133:5438/simplebank?sslmode=disable"
}

var testStore Store

func TestMain(m *testing.M) {
	// pgx.Connect - create only on connection
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
