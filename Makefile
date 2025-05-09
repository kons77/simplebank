postgres: 
	docker run --name pg-simple-bank -p 5438:5432 -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it pg-simple-bank createdb --username=postgres --owner=postgres simplebank 

dropdb:
	docker exec -it pg-simple-bank dropdb --username=postgres simplebank 

migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5438/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5438/simplebank?sslmode=disable" -verbose down

migrateuplocally:
	./migrate -path db/migration -database "postgresql://postgres:secret@192.168.88.133:5438/simplebank?sslmode=disable" -verbose up

migratedownlocally:
	./migrate -path db/migration -database "postgresql://postgres:secret@192.168.88.133:5438/simplebank?sslmode=disable" -verbose down

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

.PHONY: postgres createdb dropdb migrateup migratedown migrateuplocally migratedownlocally sqlc test server testall mock 
