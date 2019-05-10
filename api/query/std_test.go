package query_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"
)

func TestGetColumnsDefault(t *testing.T) {
  apitest.RunQueryTest(t, "site-get-columns", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/", nil)
  })
}

func TestGetColumns(t *testing.T) {
  apitest.RunQueryTest(t, "site-get-columns", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/column/", nil)
  })
}

func TestPostColumns(t *testing.T) {
  apitest.RunQueryTest(t, "site-get-columns", func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/column/", nil)
  })
}

func TestGetRows(t *testing.T) {
  apitest.RunQueryTest(t, "site-get-rows", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/row/", nil)
  })
}

func TestPostRowsDefault(t *testing.T) {
  payloadfile, err := os.Open("testdata/site-get-rows-where.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  // Test the path as the default POST operation.
  apitest.RunQueryTest(t, "site-get-rows-where", func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/", payloadfile)
  })
}

func TestPostRows(t *testing.T) {
  payloadfile, err := os.Open("testdata/site-get-rows-where.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  // Test the path as the default POST operation.
  apitest.RunQueryTest(t, "site-get-rows-where", func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/row/", payloadfile)
  })
}

func TestGetSummaries(t *testing.T) {
  apitest.RunQueryTest(t, "site-get-summary", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/summary/", nil)
  })
}
