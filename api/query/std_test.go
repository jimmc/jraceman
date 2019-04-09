package query_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"
)

func TestGetColumns(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/query/site/", nil)
  }
  if err := apitest.StartQueryToGolden("site-get-columns", request);
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

  request := func() (*http.Request, error) {
    return http.NewRequest("POST", "/api/query/site/", payloadfile)
  }
  if err := apitest.SetupToGolden(repos, handler, "site-get-rows-where", request);
       err != nil {
    t.Error(err.Error())
  }
}
