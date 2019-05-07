package crud_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"
)

// TODO: Update, Delete, Patch

func TestList(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/", nil)
  }
  r := apitest.NewCrudTester()
  r.Run(t, "site-list", request)
}

func TestListLimit(t *testing.T) {
  request1 := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=1&offset=1", nil)
  }
  r := apitest.NewCrudTester()
  r.Init(t)
  r.RunTest(t, "site-list-limit-1", request1)

  request2 := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=2&offset=2", nil)
  }
  r.RunTest(t, "site-list-limit-2", request2)
  r.Close()
}

func TestGet(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S2", nil)
  }
  r := apitest.NewCrudTester()
  r.Run(t, "site-get", request)
}

func TestCreateWithID(t *testing.T) {
  payloadfile, err := os.Open("testdata/site-create-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  request := func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/crud/site/", payloadfile)
  }
  r := apitest.NewCrudTester()
  r.Init(t)
  r.RunTest(t, "site-create-id", request)

  requestGet := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S10", nil)
  }
  r.RunTest(t, "site-create-id-get", requestGet)

  r.Close()
}
