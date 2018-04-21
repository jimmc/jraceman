package api

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jracemango/api/crud"
  "github.com/jimmc/jracemango/domain"
)

type handler struct {
  config *Config
}

type Config struct {
  Prefix string
  DomainRepos domain.Repos
}

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

  mux.HandleFunc(h.apiPrefix("foo"), h.foo)
  return mux
}

func (h *handler) foo(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Foo not implemented", http.StatusForbidden)
}

func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
