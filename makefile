createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down	

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

start: 
	docker exec -it postgres psql -U root -d simple_bank

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/hiiamanop/simple_bank/db/sqlc Store 

.PHONY: createdb dropdb migrateup migratedown sqlc test start server mock