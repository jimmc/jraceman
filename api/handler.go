package api

import (
  "fmt"
  "net/http"
)

type handler struct {
  config *Config
}

type Config struct {
  Prefix string
}

func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()
  mux.HandleFunc(h.apiPrefix("foo"), h.foo)
  // mux.HandleFunc(h.apiPrefix("bar"), h.bar)
  return mux
}

func (h *handler) foo(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Foo not implemented", http.StatusForbidden)
}

func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
