package handlers

import (
  "errors"
  "net/url"
)

func (r ShortLinkRequest) Validate(cfg configger) error {
  u, err := url.ParseRequestURI(r.Link)
  if err != nil {
    return err
  }
  if u == nil || u.Host == cfg.GetHost() {
    return errors.New("can't do that")
  }
  return nil
}
