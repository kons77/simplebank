# db source for local dev linux pc 
DSN_LINUX=postgresql://postgres:secret@192.168.88.133:5438/simplebank?sslmode=disable
# db source for github actions 
DSN_GH=postgresql://postgres:secret@localhost:5438/simplebank?sslmode=disable

postgres: 
	docker run --name pg-simple-bank -p 5438:5432 -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it pg-simple-bank createdb --username=postgres --owner=postgres simplebank 

dropdb:
	docker exec -it pg-simple-bank dropdb --username=postgres simplebank 

migrateup:
	migrate -path db/migration -database "${DSN_GH}" -verbose up

migrateup1:
	migrate -path db/migration -database "${DSN_GH}" -verbose up 1

migratedown:
	migrate -path db/migration -database "${DSN_GH}" -verbose down

migratedown1: 
	migrate -path db/migration -database "${DSN_GH}" -verbose down 1

migrateup_lcl:
	migrate -path db/migration -database "${DSN_LINUX}" -verbose up

migratedown_lcl:
	migrate -path db/migration -database "${DSN_LINUX}" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

testall: 
	@rm -f coverage.out
	@echo running all tests with covers
	@go test -v ./... -coverprofile="coverage.out" || echo "test fails"
	@go tool cover -html="coverage.out"

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/kons77/simplebank/db/sqlc Store 

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 migrateup_lcl migratedown_lcl sqlc test server testall mock 
