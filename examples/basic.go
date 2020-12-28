package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jonahgeorge/sqlexplorer"
	_ "github.com/lib/pq"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open(sqlexplorer.Postgres, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	explorer, _ := sqlexplorer.New(sqlexplorer.Config{
		Driver: sqlexplorer.Postgres,
		DB:     db,
	})

	http.ListenAndServe(":8080", explorer)
}
