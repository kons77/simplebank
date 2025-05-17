# Use DB_SOURCE from environment if set or default to localhost for safety
DB_SOURCE ?= postgresql://postgres:secret@localhost:5438/simplebank?sslmode=disable

postgres: 
	docker run --name pg-simple-bank -p 5438:5432 -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it pg-simple-bank createdb --username=postgres --owner=postgres simplebank 

dropdb:
	docker exec -it pg-simple-bank dropdb --username=postgres simplebank 

migrateup:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose up

migrateup1:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose up 1

migratedown:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose down

migratedown1: 
	migrate -path db/migration -database "${DB_SOURCE}" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

testall: 
	@echo running all tests with covers
	@go test -v ./... -coverprofile="coverage.out" || echo "test fails"
	@go tool cover -html="coverage.out"

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/kons77/simplebank/db/sqlc Store 

# run github actions locally using act, docker required
actcheck:
	act -s DB_SOURCE=${DB_SOURCE} -s TOKEN_SYMMETRIC_KEY=12345678123456781234567812345678

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 migrateup_lcl migratedown_lcl sqlc test server testall mock actcheck