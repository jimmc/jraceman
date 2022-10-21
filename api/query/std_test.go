package query_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jraceman/api/test"

  goldenbase "github.com/jimmc/golden/base"
)

func TestGetColumnsDefault(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-columns", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/", nil)
  }), "RunQueryTest")
}

func TestGetColumns(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-columns", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/column/", nil)
  }), "RunQueryTest")
}

func TestPostColumns(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-columns", func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/column/", nil)
  }), "RunQueryTest")
}

func TestGetRows(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-rows", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/row/", nil)
  }), "RunQueryTest")
}

func TestPostRowsDefault(t *testing.T) {
  payloadfile, err := os.Open("testdata/site-get-rows-where.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  // Test the path as the default POST operation.
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-rows-where", func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/", payloadfile)
  }), "RunQueryTest")
}

func TestPostRows(t *testing.T) {
  payloadfile, err := os.Open("testdata/site-get-rows-where.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  // Test the path as the default POST operation.
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-rows-where", func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/row/", payloadfile)
  }), "RunQueryTest")
}

func TestGetSummaries(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-summary", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/summary/", nil)
  }), "RunQueryTest")
}
