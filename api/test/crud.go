package test

import (
  "net/http"
  "testing"

  "github.com/jimmc/jracemango/api/crud"
)

// CreateCrudHandler create an http handler for our crud calls.
func CreateCrudHandler(r *Tester) http.Handler {
  config := &crud.Config{
    Prefix: "/api/crud/",
    DomainRepos: r.repos,
  }
  return crud.NewHandler(config)
}

// NewCrudTester returns a new Tester for testing our crud calls.
func NewCrudTester() *Tester {
  return NewTester(CreateCrudHandler)
}

// RunCrudTest creates a new Tester and runs a test for a crud call.
func RunCrudTest(t *testing.T, basename string, callback func() (*http.Request, error)) {
  t.Helper()
  r := NewCrudTester()
  r.Run(t, basename, callback)
}
