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
  if err := apitest.StartReportToGoldenWithRoots("list-get", reportRoots, request); err != nil {
    t.Error(err.Error())
  }
}
