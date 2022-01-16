package internal

import (
  "database/sql"
  "log"
)

type Config struct {
  hostname string
}

type Env struct {
  DB *sql.DB
  Cfg Config
}

func NewEnv() Env {
  db, err := sql.Open("postgres", "host=db port=5432 user=app password=password dbname=shortlinks sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }

  return Env{
    DB:  db,
    Cfg: newConfig("localhost:8080"),
  }
}

func (c Config) GetHost() string {
  return c.hostname
}

func newConfig(hostname string) Config {
  return Config{hostname: hostname}
}
