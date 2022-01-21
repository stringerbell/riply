package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"range/internal/storage"
)

type shortLinkStorage interface {
	GetOneShortLink(shortLink string) (storage.ShortLinkResponse, error)
	RecordHit(r *http.Request, resp storage.ShortLinkResponse) error
}

func NewShortLinkForwarder(sto shortLinkStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		resp, err := sto.GetOneShortLink(v["shortlink"])
		if err != nil {
			var nf storage.ErrNotFound
			if errors.As(err, &nf) {
				newError(w, nf, http.StatusNotFound)
				return
			}
			log.Println("NewShortLinkForwarder sto.GetOneShortLink err: %w", err)
			newError(w, nil, http.StatusInternalServerError)
			return
		}
		go func() {
			err = sto.RecordHit(r, resp)
			if err != nil {
				log.Println("NewShortLinkForwarder sto.RecordHit err: %w", err)
			}
		}()
		http.Redirect(w, r, resp.Link, http.StatusFound)
	}
}
