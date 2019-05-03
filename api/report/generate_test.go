package report_test

import (
  "net/http"
  "os"
  "strings"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"

  goldenhttp "github.com/jimmc/golden/http"
)

var testRoots = []string{"testdata", "../../report/template"}

func TestGet(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/generate/?name=site-report&data=S2", nil)
  }
  if err := apitest.StartReportToGoldenWithRoots("site-report", testRoots, request); err != nil {
    t.Error(err.Error())
  }
}

func TestPost(t *testing.T) {
  repos, handler, err := apitest.StartReportToSetupWithRoots(testRoots)
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  payloadfile, err := os.Open("testdata/site-report.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  request := func() (*http.Request, error) {
    req, err := http.NewRequest("POST", "/api/report/site-report/generate/", payloadfile)
    if err != nil {
      return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  }
  if err := goldenhttp.SetupToGolden(repos.DB(), handler, "site-report", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestOrderbyName(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/generate/?name=site-all-report&orderby=name", nil)
  }
  if err := apitest.StartReportDbToGoldenWithRoots("site-report", "site-orderby-name", testRoots, request); err != nil {
    t.Error(err.Error())
  }
}

/* TODO: need to define default order-by as "name".
func TestOrderbyNone(t *testing.T) {
  // The default sort for site-all-report is name, so leaving it off is like specifying "name".
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/generate/?name=site-all-report", nil)
  }
  if err := apitest.StartReportDbToGoldenWithRoots("site-report", "site-orderby-name", testRoots, request); err != nil {
    t.Error(err.Error())
  }
}
*/

func TestOrderbyCity(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/generate/?name=site-all-report&orderby=city", nil)
  }
  if err := apitest.StartReportDbToGoldenWithRoots("site-report", "site-orderby-city", testRoots, request); err != nil {
    t.Error(err.Error())
  }
}

func TestOrderbyInvalid(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/generate/?name=site-all-report&orderby=invalid", nil)
  }
  err := apitest.StartReportDbToGoldenWithRoots("site-report", "site-orderby-city", testRoots, request)
  if err == nil {
    t.Fatalf("Expected error for invalid sort key")
  }
  if !strings.Contains(err.Error(), "invalid orderby") {
    t.Errorf("Expected error about invalid orderby, but got something else: %v", err)
  }
}

func TestWherePost(t *testing.T) {
  repos, handler, err := apitest.StartReportToSetupWithRoots(testRoots)
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  payloadfile, err := os.Open("testdata/site-report-where.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  request := func() (*http.Request, error) {
    req, err := http.NewRequest("POST", "/api/report/site-report-where/generate/", payloadfile)
    if err != nil {
      return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  }
  if err := goldenhttp.SetupDbToGolden(repos.DB(), handler, "site-report", "site-report-where", request);
       err != nil {
    t.Error(err.Error())
  }
}
