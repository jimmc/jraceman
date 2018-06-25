package crud

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type handler struct {
  config *Config
}

// Config provides configuration of the http handler for CRUD calls.
type Config struct {
  Prefix string
  DomainRepos domain.Repos
}

// NewHandler creates the http handler for CRUD calls.
func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()
  mux.HandleFunc(h.crudPrefix("area"), h.area)
  mux.HandleFunc(h.crudPrefix("competition"), h.competition)
  mux.HandleFunc(h.crudPrefix("gender"), h.gender)
  mux.HandleFunc(h.crudPrefix("level"), h.level)
  mux.HandleFunc(h.crudPrefix("site"), h.site)
  return mux
}

func (h *handler) crudPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
