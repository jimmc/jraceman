package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/query"
)

// CreateQueryHandler create an http handler for our query calls.
func CreateQueryHandler(r *Tester) http.Handler {
  config := &query.Config{
    Prefix: "/api/query/",
    DomainRepos: r.repos,
  }
  return query.NewHandler(config)
}

// NewQueryTester returns a new Tester for testing our query calls.
func NewQueryTester() *Tester {
  return NewTester(CreateQueryHandler)
}

// RunQueryTest creates a new Tester and runs a test for a query call.
func RunQueryTest(basename string, callback func() (*http.Request, error)) error {
  r := NewQueryTester()
  return RunOneWith(r, basename, callback)
}
