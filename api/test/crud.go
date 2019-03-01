package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/crud"
  "github.com/jimmc/jracemango/dbrepo"
  dbtest "github.com/jimmc/jracemango/dbrepo/test"
)

func StartCrudToGolden(basename string, callback func() (*http.Request, error)) error {
  repos, handler, err := StartCrudToSetup()
  if err != nil{
    return err
  }
  defer repos.Close()

  return SetupToGolden(repos, handler, basename, callback)
}

// StartCrudToSetup initializes the database and the http handler for api/crud.
func StartCrudToSetup() (*dbrepo.Repos, http.Handler, error) {
  repos, err := dbtest.ReposEmpty()
  if err != nil {
    return nil, nil, err
  }
  config := &crud.Config{
    Prefix: "/api/crud/",
    DomainRepos: repos,
  }
  handler := crud.NewHandler(config)
  return repos, handler, nil
}
