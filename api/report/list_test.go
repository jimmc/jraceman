package report_test

import (
  "net/http"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"

  goldenbase "github.com/jimmc/golden/base"
)

func TestListGet(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/", nil)
  }
  reportRoots := []string{"testdata", "../../report/template"}
  goldenbase.FatalIfError(t, apitest.RunReportTest("list-get", reportRoots, request), "RunReportTest")
}
