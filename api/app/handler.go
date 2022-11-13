// Package app provides the api interface to app-specific functions.
package app

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jraceman/domain"

  "github.com/jimmc/auth/auth"
  "github.com/jimmc/auth/permissions"
)

type handler struct {
  config *Config
}

const (
  editRegatta = permissions.Permission("edit_regatta")
)

// Config provides configuration of the http handler for our calls.
type Config struct {
  Prefix string
  DomainRepos domain.Repos
  AuthHandler *auth.Handler
}

// NewHandler creates the http handler for our calls.
func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()
  mux.HandleFunc(h.apiPrefix("event"), c.AuthHandler.RequirePermissionFunc(h.event, editRegatta))
  return mux
}

func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
