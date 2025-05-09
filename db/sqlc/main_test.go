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

var acc11, acc12 Account
var argTestTransferParams CreateTransferParams

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

	/* сreate accounts once
	initTestAccounts()

	// сleanup after all tests (will run even if there's a panic)
	defer cleanupTestAccounts()*/

	os.Exit(m.Run())
}

/*
func initTestAccounts() (err error) {

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc11, err := testStore.CreateAccount(context.Background(), arg)
	if err != nil {
		return err
	}

	arg = CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc12, err := testStore.CreateAccount(context.Background(), arg)
	if err != nil {
		return err
	}

	argTestTransferParams = CreateTransferParams{
		FromAccountID: util.ToPgInt8(acc11.ID),
		ToAccountID:   util.ToPgInt8(acc12.ID),
		Amount:        util.RandomMoney(),
	}

	return nil
}

// cleanupTestAccounts delete accounts up after all tests (will run even if there's a panic)
func cleanupTestAccounts() {
	_ = testStore.DeleteAccount(context.Background(), acc11.ID)
	_ = testStore.DeleteAccount(context.Background(), acc12.ID)
}
*/
