package crud_test

import (
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "os"
  "testing"

  "github.com/jimmc/jracemango/api/crud"
  "github.com/jimmc/jracemango/dbrepo/dbtesting"
)

// List, Get, Create, Update, Delete, Patch

func setupToGolden(basename string, callback func() (*http.Request, error)) error {
  setupfilename := "testdata/" + basename + ".setup"
  outfilename := "testdata/" + basename + ".out"
  goldenfilename := "testdata/" + basename + ".golden"
  repos,err := dbtesting.ReposWithSetupFile(setupfilename)
  if err!= nil {
    return err
  }
  defer repos.Close()
  config := &crud.Config{
    Prefix: "/api/crud/",
    DomainRepos: repos,
  }
  handler := crud.NewHandler(config)

  req, err := callback()
  if err != nil {
    return err
  }

  rr := httptest.NewRecorder()
  handler.ServeHTTP(rr, req)

  if got, want := rr.Code, http.StatusOK; got != want {
    return fmt.Errorf("response status: got %d, want %d", got, want)
  }

  body := rr.Body.Bytes()
  if len(body) == 0 {
    return errors.New("response body should not be empty")
  }

  os.Remove(outfilename)
  if err := ioutil.WriteFile(outfilename, body, 0644); err != nil {
    return err
  }

  if err := dbtesting.CompareOutToGolden(outfilename, goldenfilename); err != nil {
    return err
  }
  return nil
}

func TestList(t *testing.T) {
  if err := setupToGolden("site-list",
      func () (*http.Request, error) { return http.NewRequest("GET", "/api/crud/site/", nil) }); err != nil {
    t.Error(err.Error())
  }
}
