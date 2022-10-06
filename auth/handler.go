package auth

import (
  authlib "github.com/jimmc/auth/auth"
  "github.com/jimmc/auth/store"
)

// NewHandler returns our auth handler, which in turn wraps other
// handlers when auth is required.
func NewHandler(maxClockSkewSeconds int) authlib.Handler {
  authStore := store.NewPwFile("/tmp/jracemagopw.txt")  // TODO - use our DB
  authHandler := authlib.NewHandler(&authlib.Config{
    Prefix: "/auth/",
    Store: authStore,
    TokenCookieName: "JRACEMAN_TOKEN",
    MaxClockSkewSeconds: maxClockSkewSeconds,
  })
  return authHandler
}
