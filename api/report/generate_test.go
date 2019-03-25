package report_test

import (
  "net/http"
  "os"
  "strings"
  "testing"

  apireport "github.com/jimmc/jracemango/api/report"
  apitest "github.com/jimmc/jracemango/api/test"
  reportmain "github.com/jimmc/jracemango/report"
  "github.com/google/go-cmp/cmp"
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
    req, err := http.NewRequest("POST", "/api/report/generate/", payloadfile)
    if err != nil {
      return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  }
  if err := apitest.SetupToGolden(repos, handler, "site-report", request);
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

func TestOrderbyNone(t *testing.T) {
  // The default sort for site-all-report is name, so leaving it off is like specifying "name".
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/report/generate/?name=site-all-report", nil)
  }
  if err := apitest.StartReportDbToGoldenWithRoots("site-report", "site-orderby-name", testRoots, request); err != nil {
    t.Error(err.Error())
  }
}

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

func TestOptionsFromParameters(t *testing.T) {
  // TODO
}

func TestWhereMapFromParameters(t *testing.T) {
  tests := []struct{
    name string
    where interface{}
    expect map[string]reportmain.WhereValue
    expectError bool
  } {
    {
      name: "no_where",
      where: nil,
      expect: map[string]reportmain.WhereValue{},
    },
    {
      name: "not_map",
      where: 123,
      expectError: true,
    },
    {
      name: "not_map_of_maps",
      where: map[string]interface{}{
        "a": 456,
      },
      expectError: true,
    },
    {
      name: "one_field",
      where: map[string]interface{}{
        "a": map[string]interface{}{
          "op": "eq",
          "value": "xyz",
        },
      },
      expect: map[string]reportmain.WhereValue{
        "a": reportmain.WhereValue{Op: "eq", Value: "xyz"},
      },
    },
    {
      name: "op_missing",
      where: map[string]interface{}{
        "a": map[string]interface{}{
          "value": "xyz",
        },
      },
      expectError: true,
    },
    {
      name: "op_not_string",
      where: map[string]interface{}{
        "a": map[string]interface{}{
          "op": 789,
        },
      },
      expectError: true,
    },
    {
      name: "no_value",
      where: map[string]interface{}{
        "a": map[string]interface{}{
          "Op": "ne",
        },
      },
      expectError: true,
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got, err := apireport.WhereMapFromParametersForTesting(tc.where)
      if tc.expectError {
        if err == nil {
          t.Fatalf("whereMapFromParameters: expected error but did not get one")
        }
      } else if err != nil {
        t.Fatalf("whereMapFromParameters: unexpected error: %v", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got); diff != "" {
          t.Errorf("where mismatch (-want +got):\n%s", diff)
        }
      }
    })
  }
}
