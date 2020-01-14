package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/debug"
)

// CreateDebugHandler create an http handler for our debug calls.
func CreateDebugHandler(r *Tester) http.Handler {
  config := &debug.Config{
    Prefix: "/api/debug/",
    DomainRepos: r.repos,
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
