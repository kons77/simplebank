package gapi

import (
	"fmt"

	db "github.com/kons77/simplebank/db/sqlc"
	"github.com/kons77/simplebank/pb"
	"github.com/kons77/simplebank/token"
	"github.com/kons77/simplebank/util"
)

// Server servers gRPC requestst for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey) // for JWT tokens
	tokenMaker, err := token.NewPasetoMaker([]byte(config.TokenSymmetricKey), []byte("some_implicit"))
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
