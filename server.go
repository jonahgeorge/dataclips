package dataclips

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
	db               *sqlx.DB
	PlaceholderQuery string
}

type Config struct {
	DB               *sql.DB
	Driver           Driver
	PlaceholderQuery string
}

func New(c Config) (*Server, error) {
	return &Server{
		db:               sqlx.NewDb(c.DB, string(c.Driver)),
		PlaceholderQuery: c.PlaceholderQuery,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := http.NewServeMux()
	m.HandleFunc("/ui", s.uiHandler)
	m.HandleFunc("/query", s.queryHandler)
	m.ServeHTTP(w, r)
}
