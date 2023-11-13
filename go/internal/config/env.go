package config

import (
	"os"

	goEnv "github.com/Netflix/go-env"
	"github.com/uulwake/godotenvsafe"
)

func NewEnv() (*Env, error) {
	var envFile string
	if os.Getenv("GO_ENV") == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env"
	}

	err := godotenvsafe.Load(envFile)
	if err != nil {
		return nil, err
	}

	env := &Env{}
	_, err = goEnv.UnmarshalFromEnviron(env)
	if err != nil {
		return nil, err
	}

	return env, nil
}

type Env struct {
	PgUrl string `env:"PG_URL"`
	EsUrl string  `env:"ES_URL"`	
	Port string  `env:"PORT"`
	Salt int `env:"SALT"`
	JwtSecret string `env:"JWT_SECRET"`
	JwtExpired string `env:"JWT_EXPIRED"`
}