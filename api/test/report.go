package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/report"
  "github.com/jimmc/jracemango/dbrepo"
  dbtest "github.com/jimmc/jracemango/dbrepo/test"
)

func StartReportToGolden(basename string, callback func() (*http.Request, error)) error {
  repos, handler, err := StartReportToSetup()
  if err != nil{
    return err
  }
  defer repos.Close()

  return SetupToGolden(repos, handler, basename, callback)
}

// StartReportToSetup initializes the database and the http handler for api/report.
func StartReportToSetup() (*dbrepo.Repos, http.Handler, error) {
  repos, err := dbtest.ReposEmpty()
  if err != nil {
    return nil, nil, err
  }
  config := &report.Config{
    Prefix: "/api/report/",
    DomainRepos: repos,
    ReportRoot: "testdata",
  }
  handler := report.NewHandler(config)
  return repos, handler, nil
}
