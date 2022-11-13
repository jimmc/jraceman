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
    r, err := http.NewRequest("GET", "/api/query/site/", nil)
    if err != nil {
      return nil, err
    }
    return apitest.AddTestUser(r, "view_venue"), nil
  }), "RunQueryTest")
}

func TestGetColumns(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-columns", func() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/api/query/site/column/", nil)
    if err != nil {
      return nil, err
    }
    return apitest.AddTestUser(r, "view_venue"), nil
  }), "RunQueryTest")
}

func TestPostColumns(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-columns", func() (*http.Request, error) {
    r, err := http.NewRequest("POST", "/api/query/site/column/", nil)
    if err != nil {
      return nil, err
    }
    return apitest.AddTestUser(r, "view_venue"), nil
  }), "RunQueryTest")
}

func TestGetRows(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-rows", func() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/api/query/site/row/", nil)
    if err != nil {
      return nil, err
    }
    return apitest.AddTestUser(r, "view_venue"), nil
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
    r, err := http.NewRequest("POST", "/api/query/site/", payloadfile)
    if err != nil {
      return nil, err
    }
    return apitest.AddTestUser(r, "view_venue"), nil
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
    r, err := http.NewRequest("POST", "/api/query/site/row/", payloadfile)
    if err != nil {
      return nil, err
    }
    return apitest.AddTestUser(r, "view_venue"), nil
  }), "RunQueryTest")
}

func TestGetSummaries(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunQueryTest("site-get-summary", func() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/api/query/site/summary/", nil)
    if err != nil {
      return nil, err
    }
    return apitest.AddTestUser(r, "view_venue"), nil
  }), "RunQueryTest")
}
