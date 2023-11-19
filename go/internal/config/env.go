package config

import (
	"os"

	goEnv "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

func NewEnv() (*Env, error) {
	var envFile string
	if os.Getenv("GO_ENV") == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env"
	}

	err := godotenv.Load(envFile)
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
	PgUrl string `env:"PG_URL,required=true"`
	EsUrl string  `env:"ES_URL,required=true"`	
	Port string  `env:"PORT,required=true"`
	Salt int `env:"SALT,,required=true"`
	JwtSecret string `env:"JWT_SECRET,required=true"`
	JwtExpired int `env:"JWT_EXPIRED,required=true"`
}