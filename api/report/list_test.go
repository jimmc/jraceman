package report_test

import (
  "net/http"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"
)

func TestListGet(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/", nil)
  }
  reportRoots := []string{"testdata", "../../report/template"}
  apitest.RunReportTest(t, "list-get", reportRoots, request)
}
