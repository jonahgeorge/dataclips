package main

import (
	dsql "database/sql"
	"log"
	"net/http"
	"os"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jonahgeorge/dataclips"
	_ "github.com/lib/pq"
)

func main() {
	go runDB()

	databaseURL := "root@tcp(localhost:3306)/db"

	db, err := dsql.Open(dataclips.MySQL, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	explorer, _ := dataclips.New(dataclips.Config{
		Driver:           dataclips.MySQL,
		DB:               db,
		PlaceholderQuery: "select * from users",
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on %s", port)

	err = http.ListenAndServe(":"+port, explorer)
	if err != nil {
		log.Fatal(err)
	}
}

func runDB() {
	engine := sqle.NewDefault()
	engine.AddDatabase(createTestDatabase())
	engine.AddDatabase(information_schema.NewInformationSchemaDatabase(engine.Catalog))

	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
		Auth:     auth.NewNativeSingle("root", "", auth.AllPermissions),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	s.Start()
}

func createTestDatabase() *memory.Database {
	const (
		dbName    = "db"
		tableName = "users"
	)

	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: sql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: sql.Timestamp, Nullable: false, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()
	table.Insert(ctx, sql.NewRow("John Doe", "john@doe.com", []string{"555-555-555"}, time.Now()))
	table.Insert(ctx, sql.NewRow("John Doe", "johnalt@doe.com", []string{}, time.Now()))
	table.Insert(ctx, sql.NewRow("Jane Doe", "jane@doe.com", []string{}, time.Now()))
	table.Insert(ctx, sql.NewRow("Evil Bob", "evilbob@gmail.com", []string{"555-666-555", "666-666-666"}, time.Now()))
	return db
}
