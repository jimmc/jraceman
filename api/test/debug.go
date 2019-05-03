package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/debug"
  "github.com/jimmc/jracemango/dbrepo"
  dbtest "github.com/jimmc/jracemango/dbrepo/test"

  goldenhttp "github.com/jimmc/golden/http"
)

func StartDebugToGolden(basename string, callback func() (*http.Request, error)) error {
  repos, handler, err := StartDebugToSetup()
  if err != nil{
    return err
  }
  defer repos.Close()

  return goldenhttp.SetupToGolden(repos.DB(), handler, basename, callback)
}

// StartDebugToSetup initializes the database and the http handler for api/debug.
func StartDebugToSetup() (*dbrepo.Repos, http.Handler, error) {
  repos, err := dbtest.ReposEmpty()
  if err != nil {
    return nil, nil, err
  }
  config := &debug.Config{
    Prefix: "/api/debug/",
    DomainRepos: repos,
  }
  handler := debug.NewHandler(config)
  return repos, handler, nil
}
