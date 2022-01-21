package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"range/internal"
	"range/internal/router"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.TimeoutHandler(router.New(internal.NewEnv()), 1*time.Second, "timeout"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
