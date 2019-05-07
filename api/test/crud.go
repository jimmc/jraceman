package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/crud"
)

func CreateCrudHandler(r *Tester) http.Handler {
  config := &crud.Config{
    Prefix: "/api/crud/",
    DomainRepos: r.repos,
  }
  return crud.NewHandler(config)
}

func NewCrudTester(basename string, callback func() (*http.Request, error)) *Tester {
  return NewTester(basename, CreateCrudHandler, callback)
}
