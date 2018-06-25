package query

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type handler struct {
  config *Config
}

// Config provides configuration of the http handler for query calls.
type Config struct {
  Prefix string
  DomainRepos domain.Repos
}

// NewHandler creates the http handler for query calls.
func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()
  mux.HandleFunc(h.queryPrefix("area"), h.area)
  mux.HandleFunc(h.queryPrefix("competition"), h.competition)
  mux.HandleFunc(h.queryPrefix("gender"), h.gender)
  mux.HandleFunc(h.queryPrefix("level"), h.level)
  mux.HandleFunc(h.queryPrefix("site"), h.site)
  return mux
}

func (h *handler) queryPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
