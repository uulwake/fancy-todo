package database

import (
	"database/sql"
	"fancy-todo/internal/config"

	_ "github.com/lib/pq"
)

func NewPg(env *config.Env) (*sql.DB, error) {
	pg, err := sql.Open("postgres", env.PgUrl)
	if err != nil {
		return nil, err
	}

	pg.SetMaxOpenConns(10)
	pg.SetMaxIdleConns(2)

	return pg, nil;
}