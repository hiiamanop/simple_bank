package main

import (
	"database/sql"
	"log"

	"github.com/hiiamanop/simple_bank/api"
	db "github.com/hiiamanop/simple_bank/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	dbConn, err := sql.Open("postgres", "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(dbConn)
	server := api.NewServer(store)

	err = server.Start(":8080")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
