postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:17

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres17 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/zhifaq/simple_bank/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock