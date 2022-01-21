package handlers

import (
  "encoding/json"
  "errors"
  "github.com/gorilla/mux"
  "net/http"
  "range/internal/storage"
)

type shortLinkStatsStorage interface {
  GetStatsForShortLink(shortlink string) (storage.ShortLinkStats, error)
}

func NewStatsHandler(sto shortLinkStatsStorage) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    v := mux.Vars(r)
    stats, err := sto.GetStatsForShortLink(v["shortlink"])
    if err != nil {
      var nf storage.ErrNotFound
      if errors.As(err, &nf) {
        newError(w, nf, http.StatusNotFound)
        return
      }
    }
    bs, _ := json.Marshal(stats)
    w.Header().Set("Content-Type", "application/json")
    w.Write(bs)
  }
}
