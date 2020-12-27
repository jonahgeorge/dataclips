package sqlexplorer

import (
	"database/sql"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Driver string

const (
	MySQL    = "mysql"
	Postgres = "postgres"
)

type Server struct {
	db *sqlx.DB
}

type Config struct {
	DB     *sql.DB
	Driver Driver
}

func New(c Config) (*Server, error) {
	return &Server{db: sqlx.NewDb(c.DB, string(c.Driver))}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := http.NewServeMux()
	m.HandleFunc("/ui", s.uiHandler)
	m.HandleFunc("/query", s.queryHandler)
	m.ServeHTTP(w, r)
}
