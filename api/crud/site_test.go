package crud_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"
)

// TODO: Update, Delete, Patch

func TestList(t *testing.T) {
  apitest.RunCrudTest(t, "site-list", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/", nil)
  })
}

func TestListLimit(t *testing.T) {
  r := apitest.NewCrudTester()
  r.Init(t)

  r.RunTest(t, "site-list-limit-1", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=1&offset=1", nil)
  })

  r.RunTest(t, "site-list-limit-2", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=2&offset=2", nil)
  })

  r.Close()
}

func TestGet(t *testing.T) {
  apitest.RunCrudTest(t, "site-get", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S2", nil)
  })
}

func TestCreateWithID(t *testing.T) {
  r := apitest.NewCrudTester()
  r.Init(t)

  payloadfile, err := os.Open("testdata/site-create-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  r.RunTest(t, "site-create-id", func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/crud/site/", payloadfile)
  })

  r.RunTest(t, "site-create-id-get", func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S10", nil)
  })

  r.Close()
}
