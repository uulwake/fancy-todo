package database

import (
	"fancy-todo/internal/config"

	"github.com/elastic/go-elasticsearch/v7"
)

func NewEs(env *config.Env) (*elasticsearch.Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: []string{
			env.EsUrl,
		},
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		 return nil, err
	}

	return es, nil
}