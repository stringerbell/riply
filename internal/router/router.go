package router

import (
  "github.com/gorilla/mux"
)


func New(env internal.Env) *mux.Router {
  r := mux.NewRouter()
  return r
}
