package crud_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"

  goldenbase "github.com/jimmc/golden/base"
  goldenhttp "github.com/jimmc/golden/http"
)

// TODO: Update, Delete, Patch

func TestList(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunCrudTest("site-list", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/", nil)
  }), "RunCrudTest")
}

func TestListLimit(t *testing.T) {
  r := apitest.NewCrudTester()
  goldenbase.FatalIfError(t, r.Init(), "Init")

  goldenbase.FatalIfError(t, goldenhttp.RunTestWith(r, "site-list-limit-1", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=1&offset=1", nil)
  }), "RunTestWith")

  goldenbase.FatalIfError(t, goldenhttp.RunTestWith(r, "site-list-limit-2", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=2&offset=2", nil)
  }), "RunTestWith")

  r.Close()
}

func TestGet(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunCrudTest("site-get", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S2", nil)
  }), "RunCrudTest")
}

func TestCreateWithID(t *testing.T) {
  r := apitest.NewCrudTester()
  goldenbase.FatalIfError(t, r.Init(), "Init")

  payloadfile, err := os.Open("testdata/site-create-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  goldenbase.FatalIfError(t, goldenhttp.RunTestWith(r, "site-create-id", func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/crud/site/", payloadfile)
  }), "RunTestWith")

  goldenbase.FatalIfError(t, goldenhttp.RunTestWith(r, "site-create-id-get", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S10", nil)
  }), "RunTestWith")

  r.Close()
}
