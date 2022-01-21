package handlers

import (
  "encoding/json"
  "errors"
  "net/http"
)

type err struct {
  Detail string `json:"detail"`
  Error bool `json:"error"`
}
func newError(w http.ResponseWriter, e error, statusCode int) []byte {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  if e == nil {
    e = errors.New("internal")
  }
  bs, _ := json.Marshal(err{
    Error: true,
    Detail: e.Error(),
  })
  return bs
}
