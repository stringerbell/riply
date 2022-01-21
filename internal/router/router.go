package router

import (
  "github.com/gorilla/mux"
  "net/http"
  "range/internal"
  "range/internal/handlers"
  "range/internal/storage"
)


func New(env internal.Env) *mux.Router {
  r := mux.NewRouter()
  sto := storage.NewShortLinkStorage(env.DB)
  r.HandleFunc("/new", handlers.NewShortLinkHandler(sto, env.Cfg)).Methods(http.MethodPost)
  r.HandleFunc("/{shortlink}+", handlers.NewStatsHandler(sto)).Methods(http.MethodGet)
  r.HandleFunc("/{shortlink}", handlers.NewShortLinkForwarder(sto)).Methods(http.MethodGet)

  return r
}
