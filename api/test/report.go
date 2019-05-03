package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/report"
  "github.com/jimmc/jracemango/dbrepo"
  dbtest "github.com/jimmc/jracemango/dbrepo/test"

  goldenhttp "github.com/jimmc/golden/http"
)

func StartReportToGolden(basename string, callback func() (*http.Request, error)) error {
  return StartReportDbToGoldenWithRoots(basename, basename, []string{"testdata"}, callback)
}

func StartReportDbToGolden(dbbasename, basename string, callback func() (*http.Request, error)) error {
  return StartReportDbToGoldenWithRoots(dbbasename, basename, []string{"testdata"}, callback)
}

func StartReportToGoldenWithRoots(basename string, roots []string, callback func() (*http.Request, error)) error {
  return StartReportDbToGoldenWithRoots(basename, basename, roots, callback)
}

func StartReportDbToGoldenWithRoots(dbbasename, basename string, roots []string, callback func() (*http.Request, error)) error {
  repos, handler, err := StartReportToSetupWithRoots(roots)
  if err != nil{
    return err
  }
  defer repos.Close()

  return goldenhttp.SetupDbToGolden(repos.DB(), handler, dbbasename, basename, callback)
}

// StartReportToSetup initializes the database and the http handler for api/report.
func StartReportToSetup() (*dbrepo.Repos, http.Handler, error) {
  return StartReportToSetupWithRoots([]string{"testdata"})
}

func StartReportToSetupWithRoots(roots []string) (*dbrepo.Repos, http.Handler, error) {
  repos, err := dbtest.ReposEmpty()
  if err != nil {
    return nil, nil, err
  }
  config := &report.Config{
    Prefix: "/api/report/",
    DomainRepos: repos,
    ReportRoots: roots,
  }
  handler := report.NewHandler(config)
  return repos, handler, nil
}
