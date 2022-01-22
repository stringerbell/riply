package internal

import (
	"database/sql"
	"log"
	"os"
)

type Config struct {
	hostname string
}

type Env struct {
	DB  *sql.DB
	Cfg Config
}

func NewEnv() Env {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return Env{
		DB:  db,
		Cfg: newConfig(os.Getenv("HOST")),
	}
}

func (c Config) GetHost() string {
	return c.hostname
}

func newConfig(hostname string) Config {
	return Config{hostname: hostname}
}
