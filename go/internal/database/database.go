package database

import (
	"database/sql"
	"fancy-todo/internal/config"

	"github.com/elastic/go-elasticsearch/v7"
)

func NewDb(env *config.Env) (*Db, error) {
	pg, err := NewPg(env)
	if err != nil {
		return nil, err
	}

	es, err := NewEs(env)
	if err != nil {
		return nil, err
	}

	db := &Db{
		Pg: pg,
		Es: es,
	}

	return db, nil
}

type Db struct {
	Pg *sql.DB
	Es *elasticsearch.Client
}