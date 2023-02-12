package auth

import (
  "github.com/jimmc/jraceman/dbrepo/compat"

  authlib "github.com/jimmc/auth/auth"
)

// NewHandler returns our auth handler, which in turn wraps other
// handlers when auth is required.
func NewHandler(db compat.DBorTx) *authlib.Handler {
  authStore := NewPwDB(db)
  authHandler := authlib.NewHandler(&authlib.Config{
    Prefix: "/auth/",
    Store: authStore,
    TokenCookieName: "JRACEMAN_TOKEN",
  })
  return authHandler
}
