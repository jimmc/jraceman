package database

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jracemango/domain"
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
  mux.HandleFunc(h.apiPrefix("upgrade"), h.upgrade)
  mux.HandleFunc(h.config.Prefix, h.blank)
  return mux
}

func (h *handler) blank(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Try /api/database/upgrade", http.StatusForbidden)
}

func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
