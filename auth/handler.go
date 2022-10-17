package auth

import (
  "database/sql"

  authlib "github.com/jimmc/auth/auth"
)

// NewHandler returns our auth handler, which in turn wraps other
// handlers when auth is required.
func NewHandler(db *sql.DB, maxClockSkewSeconds int) *authlib.Handler {
  authStore := NewPwDB(db)
  authHandler := authlib.NewHandler(&authlib.Config{
    Prefix: "/auth/",
    Store: authStore,
    TokenCookieName: "JRACEMAN_TOKEN",
    MaxClockSkewSeconds: maxClockSkewSeconds,
  })
  return authHandler
}
