package crud_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"

  goldenbase "github.com/jimmc/golden/base"
)

// TODO: Update, Delete, Patch

func TestList(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/", nil)
  }
  r := apitest.NewCrudTester("site-list", request)
  goldenbase.RunT(t, r)
}

func TestListLimit(t *testing.T) {
  request1 := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=1&offset=1", nil)
  }
  r := apitest.NewCrudTester("site-list-limit-1", request1)
  if err := r.SetupDb(); err != nil {
    t.Fatalf("Error in SetupDb: %v", err)
  }
  r.LoadActAssertT(t)

  request2 := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=2&offset=2", nil)
  }
  r.BaseName = "site-list-limit-2"
  r.Callback = request2
  r.LoadActAssertT(t)
  r.Close()
}

func TestGet(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S2", nil)
  }
  r := apitest.NewCrudTester("site-get", request)
  goldenbase.RunT(t, r)
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
  r := apitest.NewCrudTester("site-create-id", request)
  if err := r.SetupDb(); err != nil {
    t.Fatalf("Error setting up database: %v", err)
  }
  r.LoadActAssertT(t)

  requestGet := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S10", nil)
  }
  r.BaseName = "site-create-id-get"
  r.Callback = requestGet
  r.LoadActAssertT(t)

  r.Close()
}
