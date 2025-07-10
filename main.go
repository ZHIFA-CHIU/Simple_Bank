package main

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/zhifaq/simple_bank/api"
	db "github.com/zhifaq/simple_bank/sqlc"
	"github.com/zhifaq/simple_bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err := server.Start(config.ServerAddress); err != nil {
		panic(err)
	}
}
