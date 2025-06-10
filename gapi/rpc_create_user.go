package gapi

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/kons77/simplebank/db/sqlc"
	"github.com/kons77/simplebank/pb"
	"github.com/kons77/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// *pq.Error for lib/pq, *pgconn.PgError for jackc/pgx/v5

		if pqErr, ok := err.(*pgconn.PgError); ok {
			log.Println(pqErr.Code, pqErr.Message)
			switch pqErr.Code {
			case "23505": //23505 unique violation
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// we cannot just return the db.User object here, because it’s a different type.
	// we should not mix up the DB layer struct with the API layer struct. It’s better to separate them from each other,
	// because sometimes we don’t want to return every field in the DB to the client.
	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
