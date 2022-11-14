package report_test

import (
  "net/http"
  "testing"

  apitest "github.com/jimmc/jraceman/api/test"

  goldenbase "github.com/jimmc/golden/base"
)

func TestListGet(t *testing.T) {
  request := func() (*http.Request, error) {
    req, err := http.NewRequest("GET", "/api/report/", nil)
    if err != nil {
      return nil, err
    }
    req = apitest.AddTestUser(req, "view_regatta")
    return req, nil
  }
  reportRoots := []string{"testdata", "../../report/template"}
  goldenbase.FatalIfError(t, apitest.RunReportTest("list-get", reportRoots, request), "RunReportTest")
}
