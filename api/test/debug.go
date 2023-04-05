package test

import (
  "net/http"

  "github.com/jimmc/jraceman/api/debug"
  "github.com/jimmc/jraceman/auth"
  "github.com/jimmc/jraceman/dbrepo"
  dbtest "github.com/jimmc/jraceman/dbrepo/test"
)

// CreateDebugHandlerForRepos creates an http handler for our debug calls
// using the specified repos.
func CreateDebugHandlerForRepos(r *Tester, dbr *dbrepo.Repos) http.Handler {
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
func NewDebugTester(dbr *dbrepo.Repos) *Tester {
  createDebugHandler := func(r *Tester) http.Handler {
    return CreateDebugHandlerForRepos(r, dbr)
  }
  return NewTester(createDebugHandler)
}

// RunDebugTest creates a new Tester and runs a test for a debug call.
func RunDebugTest(basename string, callback func() (*http.Request, error)) error {
  dbr, cleanup, err := dbtest.ReposEmpty()
  if err != nil {
    panic("Failed to open in-memory database")
  }
  r := NewDebugTester(dbr)
  err = RunOneWith(r, basename, callback)
  cleanup()
  return err
}
