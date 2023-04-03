package test

import (
  "net/http"

  "github.com/jimmc/jraceman/api/debug"
  "github.com/jimmc/jraceman/auth"
  dbtest "github.com/jimmc/jraceman/dbrepo/test"
)

// CreateDebugHandler create an http handler for our debug calls.
func CreateDebugHandler(r *Tester) http.Handler {
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    panic("Failed to open in-memory database")
  }
  emptyDb := dbr.DB()
  authHandler := auth.NewHandler(emptyDb)
  config := &debug.Config{
    Prefix: "/api/debug/",
    DomainRepos: r.repos,
    AuthHandler: authHandler,
  }
  return debug.NewHandler(config)
}

// NewDebugTester returns a new Tester for testing our debug calls.
func NewDebugTester() *Tester {
  return NewTester(CreateDebugHandler)
}

// RunDebugTest creates a new Tester and runs a test for a debug call.
func RunDebugTest(basename string, callback func() (*http.Request, error)) error {
  r := NewDebugTester()
  return RunOneWith(r, basename, callback)
}
