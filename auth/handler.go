package auth

import (
  "github.com/jimmc/jraceman/dbrepo/conn"

  authlib "github.com/jimmc/auth/auth"
)

// NewHandler returns our auth handler, which in turn wraps other
// handlers when auth is required.
func NewHandler(db conn.DB) *authlib.Handler {
  authStore := NewPwDB(db)
  authHandler := authlib.NewHandler(&authlib.Config{
    Prefix: "/auth/",
    Store: authStore,
    TokenCookieName: "JRACEMAN_TOKEN",
  })
  return authHandler
}
