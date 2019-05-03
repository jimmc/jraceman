package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/query"
  "github.com/jimmc/jracemango/dbrepo"
  dbtest "github.com/jimmc/jracemango/dbrepo/test"

  goldenhttp "github.com/jimmc/golden/http"
)

func StartQueryToGolden(basename string, callback func() (*http.Request, error)) error {
  repos, handler, err := StartQueryToSetup()
  if err != nil{
    return err
  }
  defer repos.Close()

  return goldenhttp.SetupToGolden(repos.DB(), handler, basename, callback)
}

// StartQueryToSetup initializes the database and the http handler for api/query.
func StartQueryToSetup() (*dbrepo.Repos, http.Handler, error) {
  repos, err := dbtest.ReposEmpty()
  if err != nil {
    return nil, nil, err
  }
  config := &query.Config{
    Prefix: "/api/query/",
    DomainRepos: repos,
  }
  handler := query.NewHandler(config)
  return repos, handler, nil
}
