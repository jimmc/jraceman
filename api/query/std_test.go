package query_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"

  goldenhttp "github.com/jimmc/golden/http"
)

func TestGetColumnsDefault(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/", nil)
  }
  if err := apitest.StartQueryToGolden("site-get-columns", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestGetColumns(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/column/", nil)
  }
  if err := apitest.StartQueryToGolden("site-get-columns", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestPostColumns(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/column/", nil)
  }
  if err := apitest.StartQueryToGolden("site-get-columns", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestGetRows(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/row/", nil)
  }
  if err := apitest.StartQueryToGolden("site-get-rows", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestPostRowsDefault(t *testing.T) {
  repos, handler, err := apitest.StartQueryToSetup()
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  payloadfile, err := os.Open("testdata/site-get-rows-where.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  // Test the path as the default POST operation.
  request := func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/", payloadfile)
  }
  if err := goldenhttp.SetupToGolden(repos.DB(), handler, "site-get-rows-where", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestPostRows(t *testing.T) {
  repos, handler, err := apitest.StartQueryToSetup()
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  payloadfile, err := os.Open("testdata/site-get-rows-where.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  // Test the path as the default POST operation.
  request := func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/row/", payloadfile)
  }
  if err := goldenhttp.SetupToGolden(repos.DB(), handler, "site-get-rows-where", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestGetSummaries(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/summary/", nil)
  }
  if err := apitest.StartQueryToGolden("site-get-summary", request);
       err != nil {
    t.Error(err.Error())
  }
}

