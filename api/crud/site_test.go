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
  if err := apitest.StartCrudToGolden("site-list", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestListLimit(t *testing.T) {
  repos, handler, err := apitest.StartCrudToSetup()
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  request1 := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=1&offset=1", nil)
  }
  if err := apitest.SetupToGolden(repos, handler, "site-list-limit-1", request1); err != nil {
    t.Fatal(err.Error())
  }

  request2 := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/?limit=2&offset=2", nil)
  }
  if err := apitest.SetupToGolden(repos, handler, "site-list-limit-2", request2); err != nil {
    t.Fatal(err.Error())
  }
}

func TestGet(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S2", nil)
  }
  if err := apitest.StartCrudToGolden("site-get", request);
       err != nil {
    t.Error(err.Error())
  }
}

func TestCreateWithID(t *testing.T) {
  repos, handler, err := apitest.StartCrudToSetup()
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  payloadfile, err := os.Open("testdata/site-create-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  request := func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/crud/site/", payloadfile)
  }
  if err := apitest.SetupToGolden(repos, handler, "site-create-id", request);
       err != nil {
    t.Error(err.Error())
  }

  requestGet := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/crud/site/S10", nil)
  }
  if err := apitest.SetupToGolden(repos, handler, "site-create-id-get", requestGet);
       err != nil {
    t.Error(err.Error())
  }
}
