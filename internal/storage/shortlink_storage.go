package storage

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func NewShortLinkStorage(db *sql.DB) *ShortLinkStorage {
	return &ShortLinkStorage{
		db: db,
	}
}

type ShortLinkStorage struct {
	db *sql.DB
}

type ShortLinkRequest struct {
	Link         string
	CustomSuffix *string
}

type ShortLinkResponse struct {
	id        int64 // primary key -- keep this unexported
	Link      string
	Hash      *string   `json:"hash,omitempty"`
	Suffix    *string   `json:"suffix,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrConflict struct {
	error
	Location string
}

func (c ErrConflict) Error() string {
	if c.error == nil {
		return "<nil>"
	}
	return c.error.Error()
}

type ErrNotFound struct {
	error
}

func (c ErrNotFound) Error() string {
	return "not found"
}

func (s *ShortLinkStorage) getLinkFromRequest(url string, suffix *string) (ShortLinkResponse, error) {
	r := s.db.QueryRow("select id, link, hash, custom_suffix, created_at from shortlinks where link = $1 or custom_suffix = $2", url, ptrToString(suffix))
	var id int64
	var link string
	var hash, suf *string
	var created time.Time
	err := r.Scan(&id, &link, &hash, &suf, &created)
	if r.Err() != nil {
		return ShortLinkResponse{}, err
	}
	return ShortLinkResponse{
		id:        id,
		Link:      link,
		Hash:      hash,
		Suffix:    suf,
		CreatedAt: created,
	}, nil
}

func (s *ShortLinkStorage) GetOneShortLink(shortlink string) (ShortLinkResponse, error) {
	r := s.db.QueryRow("select id, link, hash, custom_suffix, created_at from shortlinks where hash = $1 or custom_suffix = $1", shortlink)
	var id int64
	var link string
	var hash, suff *string
	var created time.Time
	err := r.Scan(&id, &link, &hash, &suff, &created)
	if err != nil {
		if err == sql.ErrNoRows {
			return ShortLinkResponse{}, ErrNotFound{}
		}
		return ShortLinkResponse{}, fmt.Errorf("r.Scan err %w", err)
	}
	return ShortLinkResponse{
		id:        id,
		Link:      link,
		Hash:      hash,
		Suffix:    suff,
		CreatedAt: created,
	}, nil
}

func (s *ShortLinkStorage) Save(request ShortLinkRequest) (ShortLinkResponse, error) {
	link, err := s.getLinkFromRequest(request.Link, request.CustomSuffix)
	if err != nil && err != sql.ErrNoRows {
		return ShortLinkResponse{}, fmt.Errorf("s.getLinkFromRequest err: %w", err)
	}
	if !link.CreatedAt.IsZero() {
		return ShortLinkResponse{}, ErrConflict{
			error:    errors.New("that link already exists"),
			Location: location(link.Hash, link.Suffix),
		}
	}
	err = s.saveLink(request)
	if err != nil {
		return ShortLinkResponse{}, fmt.Errorf("s.saveLink err: %w", err)
	}
	return s.getLinkFromRequest(request.Link, request.CustomSuffix)
}

func location(hash *string, suffix *string) string {
	if hash != nil {
		return *hash
	}
	return *suffix
}

func ptrToString(suffix *string) string {
	if suffix == nil {
		return ""
	}
	return *suffix
}

func (s *ShortLinkStorage) saveLink(req ShortLinkRequest) error {
	if req.CustomSuffix != nil {
		_, err := s.db.Exec("insert into shortlinks(link, custom_suffix) values ($1, $2)", req.Link, req.CustomSuffix)
		return err
	}
	_, err := s.db.Exec("insert into shortlinks(link, hash) values ($1, $2)", req.Link, uniqueHash())
	return err
}

func uniqueHash() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
