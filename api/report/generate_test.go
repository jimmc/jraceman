package report_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jraceman/api/test"

  goldenbase "github.com/jimmc/golden/base"
  goldenhttp "github.com/jimmc/golden/http"
)

var testRoots = []string{"testdata", "../../report/template"}

func TestGet(t *testing.T) {
  request := func() (*http.Request, error) {
    req, err := http.NewRequest("GET", "/api/report/generate/?name=site-report&data=S2", nil)
    if err != nil {
      return nil, err
    }
    req = apitest.AddTestUser(req, "view_regatta")
    return req, nil
  }
  goldenbase.FatalIfError(t, apitest.RunReportTest("site-report", testRoots, request), "RunReportTest")
}

func TestPost(t *testing.T) {
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
    req = apitest.AddTestUser(req, "view_regatta")
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  }
  goldenbase.FatalIfError(t, apitest.RunReportTest("site-report", testRoots, request), "RunReportTest")
}

func TestOrderbyName(t *testing.T) {
  request := func() (*http.Request, error) {
    req, err := http.NewRequest("GET", "/api/report/generate/?name=site-all-report&orderby=name", nil)
    if err != nil {
      return nil, err
    }
    req = apitest.AddTestUser(req, "view_regatta")
    return req, nil
  }
  r := apitest.NewReportTester(testRoots)
  r.SetupBaseName = "site-report"
  goldenbase.FatalIfError(t, goldenhttp.RunOneWith(r, "site-orderby-name", request), "Run")
}

/* TODO: need to define default order-by as "name".
func TestOrderbyNone(t *testing.T) {
  // The default sort for site-all-report is name, so leaving it off is like specifying "name".
  request := func() (*http.Request, error) {
    req, err := http.NewRequest("GET", "/api/report/generate/?name=site-all-report", nil)
    if err != nil {
      return nil, err
    }
    req = apitest.AddTestUser(req, "view_regatta")
  }
  r := apitest.NewReportTester(testRoots)
  r.SetupBaseName = "site-report"
  goldenbase.FatalIfError(t, goldenhttp.RunOneWith(r, "site-orderby-name", request), "Run")
}
*/

func TestOrderbyCity(t *testing.T) {
  request := func() (*http.Request, error) {
    req, err := http.NewRequest("GET", "/api/report/generate/?name=site-all-report&orderby=city", nil)
    if err != nil {
      return nil, err
    }
    req = apitest.AddTestUser(req, "view_regatta")
    return req, nil
  }
  r := apitest.NewReportTester(testRoots)
  r.SetupBaseName = "site-report"
  goldenbase.FatalIfError(t, goldenhttp.RunOneWith(r, "site-orderby-city", request), "Run")
}

func TestOrderbyInvalid(t *testing.T) {
  request := func() (*http.Request, error) {
    req, err := http.NewRequest("GET", "/api/report/generate/?name=site-all-report&orderby=invalid", nil)
    if err != nil {
      return nil, err
    }
    req = apitest.AddTestUser(req, "view_regatta")
    return req, nil
  }
  r := apitest.NewReportTester(testRoots)
  r.SetupBaseName = "site-report"
  if err := r.Init(); err != nil {
    t.Fatal(err)
  }
  r.SetBaseNameAndCallback("site-orderby-city", request)
  if err := r.Arrange(); err != nil {
    t.Fatal(err)
  }
  err := r.Act()
  if err == nil {
    t.Fatalf("Expected error for invalid sort key")
  }
  r.Close()
}

func TestWherePost(t *testing.T) {
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
    req = apitest.AddTestUser(req, "view_regatta")
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  }

  r := apitest.NewReportTester(testRoots)
  r.SetupBaseName = "site-report"
  goldenbase.FatalIfError(t, goldenhttp.RunOneWith(r, "site-report-where", request), "Run")
}
