package main

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"range/internal"
	"range/internal/router"
	"time"
)

func main() {
  env := internal.NewEnv()
  defer env.DB.Close()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.TimeoutHandler(router.New(env), 1*time.Second, "timeout"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
