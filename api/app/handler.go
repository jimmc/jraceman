// Package app provides the api interface to app-specific functions.
package app

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type handler struct {
  config *Config
}

// Config provides configuration of the http handler for our calls.
type Config struct {
  Prefix string
  DomainRepos domain.Repos
}

// NewHandler creates the http handler for our calls.
func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()
  mux.HandleFunc(h.apiPrefix("event"), h.event)
  return mux
}

func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
