package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"range/internal/storage"
	"time"
)

type ShortLinkRequest struct {
	Link         string  `json:"link"`
	CustomSuffix *string `json:"custom_suffix,omitempty"`
}

type ShortLinkResponse struct {
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"created_at"`
}

type shortLinkWriter interface {
	Save(request storage.ShortLinkRequest) (storage.ShortLinkResponse, error)
}

type configger interface {
	GetHost() string
}

func NewShortLinkHandler(sto shortLinkWriter, cfg configger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Default().Println(fmt.Errorf("NewShortLinkHandler ioutil.ReadAll err: %w", err))
			newError(w, nil, http.StatusInternalServerError)
			return
		}
		req := ShortLinkRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Default().Println(fmt.Errorf("NewShortLinkHandler json.Unmarshal err: %w", err))
			w.Write(newError(w, err, http.StatusBadRequest))
			return
		}
		err = req.Validate(cfg)
		if err != nil {
			w.Write(newError(w, err, http.StatusBadRequest))
			return
		}
		res, err := sto.Save(storage.ShortLinkRequest{
			Link:         req.Link,
			CustomSuffix: req.CustomSuffix,
		})
		if err != nil {
			var conflict storage.ErrConflict
			if errors.As(err, &conflict) {
				w.Header().Set("Location", "/"+conflict.Location)
				w.Write(newError(w, err, http.StatusConflict))
				return
			}
			log.Default().Println(fmt.Errorf("NewShortLinkHandler sto.Save err: %w", err))
			w.Write(newError(w, nil, http.StatusInternalServerError))
			return
		}
		bs, err := json.Marshal(shortLinkRes(res, cfg))
		if err != nil {
			log.Default().Println(fmt.Errorf("NewShortLinkHandler json.Marshal err: %w", err))
			w.Write(newError(w, nil, http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bs)
	}
}

func shortLinkRes(res storage.ShortLinkResponse, cfg configger) ShortLinkResponse {
	return ShortLinkResponse{
		Link:      link(res, cfg),
		CreatedAt: res.CreatedAt,
	}
}

func link(res storage.ShortLinkResponse, cfg configger) string {
	l := res.Hash
	if res.Suffix != nil {
		l = res.Suffix
	}
	return fmt.Sprintf("http://%s/%s", cfg.GetHost(), *l)
}
