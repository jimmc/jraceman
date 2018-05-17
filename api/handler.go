package api

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jracemango/api/apidebug"
  "github.com/jimmc/jracemango/api/crud"
  "github.com/jimmc/jracemango/domain"
)

type handler struct {
  config *Config
}

// Config is used to configure the api handler.
// The Prefix field must be set to the part of the URL path that was
// used to route the request to this handler.
type Config struct {
  Prefix string
  DomainRepos domain.Repos
}

// NewHandler creates the http handler that is used to route api requests.
func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()

  crudPrefix := h.apiPrefix("crud")
  crudConfig := &crud.Config{
    Prefix: crudPrefix,
    DomainRepos: c.DomainRepos,
  }
  crudHandler := crud.NewHandler(crudConfig)
  mux.Handle(crudPrefix, crudHandler)

  debugPrefix := h.apiPrefix("debug")
  debugConfig := &apidebug.Config{
    Prefix: debugPrefix,
    DomainRepos: c.DomainRepos,
  }
  debugHandler := apidebug.NewHandler(debugConfig)
  mux.Handle(debugPrefix, debugHandler)

  mux.HandleFunc(h.config.Prefix, h.blank)
  return mux
}

func (h *handler) blank(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Try /api/crud or /api/debug", http.StatusForbidden)
}

// ApiPrefix composes our prefix with the next path component so that we can
// provide the right prefix to the handler that handles that next
// path component.
func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
