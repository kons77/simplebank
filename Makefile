postgres: 
	docker run --name pg-simple-bank -p 5438:5432 -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it pg-simple-bank createdb --username=postgres --owner=postgres simplebank 

dropdb:
	docker exec -it pg-simple-bank dropdb --username=postgres simplebank 

migrateup:
	./migrate -path db/migration -database "postgresql://postgres:secret@192.168.88.133:5438/simplebank?sslmode=disable" -verbose up

migratedown:
	./migrate -path db/migration -database "postgresql://postgres:secret@192.168.88.133:5438/simplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test