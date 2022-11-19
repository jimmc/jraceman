package test

import (
  "net/http"
  "net/http/httptest"

  "github.com/jimmc/jraceman/api/crud"
  "github.com/jimmc/jraceman/dbrepo"
)

// CreateCrudHandler create an http handler for our crud calls.
func CreateCrudHandler(r *Tester) http.Handler {
  return NewCrudHandler(r.repos)
}

func NewCrudHandler(repos *dbrepo.Repos) http.Handler {
  config := &crud.Config{
    Prefix: "/api/crud/",
    DomainRepos: repos,
  }
  return crud.NewHandler(config)
}

// NewCrudTester returns a new Tester for testing our crud calls.
func NewCrudTester() *Tester {
  return NewTester(CreateCrudHandler)
}

// RunCrudTest creates a new Tester and runs a test for a crud call.
func RunCrudTest(basename string, callback func() (*http.Request, error)) error {
  r := NewCrudTester()
  return RunOneWith(r, basename, callback)
}

// ServeCrudRequest create our http handler and processes the given request.
func ServeCrudRequest(req *http.Request, repos *dbrepo.Repos) *httptest.ResponseRecorder {
  handler := NewCrudHandler(repos)
  rr := httptest.NewRecorder()
  handler.ServeHTTP(rr, req)
  return rr
}
