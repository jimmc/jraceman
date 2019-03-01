package report_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"
)

func TestGet(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/generate/?name=site-report&data=S2", nil)
  }
  if err := apitest.StartReportToGolden("generate-get-site-report", request); err != nil {
    t.Error(err.Error())
  }
}

func TestPost(t *testing.T) {
  repos, handler, err := apitest.StartReportToSetup()
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  payloadfile, err := os.Open("testdata/generate-get-site-report.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  request := func() (*http.Request, error) {
    req, err := http.NewRequest("POST", "/api/report/generate/", payloadfile)
    if err != nil {
      return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  }
  if err := apitest.SetupToGolden(repos, handler, "generate-get-site-report", request);
       err != nil {
    t.Error(err.Error())
  }
}
