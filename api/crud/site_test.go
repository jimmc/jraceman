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
  "github.com/jimmc/jracemango/dbrepo"
  "github.com/jimmc/jracemango/dbrepo/dbtest"
)

// List, Get, Create, Update, Delete, Patch

func startToGolden(basename string, callback func() (*http.Request, error)) error {
  repos, handler, err := startToSetup()
  if err != nil{
    return err
  }
  defer repos.Close()

  return setupToGolden(repos, handler, basename, callback)
}

// startToSetup initializes the database and the http handler.
func startToSetup() (*dbrepo.Repos, http.Handler, error) {
  repos, err := dbtest.ReposEmpty()
  if err != nil {
    return nil, nil, err
  }
  config := &crud.Config{
    Prefix: "/api/crud/",
    DomainRepos: repos,
  }
  handler := crud.NewHandler(config)
  return repos, handler, nil
}

func setupToGolden(repos *dbrepo.Repos, handler http.Handler, basename string,
    callback func() (*http.Request, error)) error {
  setupfilename := "testdata/" + basename + ".setup"
  outfilename := "testdata/" + basename + ".out"
  goldenfilename := "testdata/" + basename + ".golden"

  if err := dbtest.LoadSetupFile(repos.DB(), setupfilename); err != nil {
    return err
  }

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

  if err := dbtest.CompareOutToGolden(outfilename, goldenfilename); err != nil {
    return err
  }
  return nil
}

func siteListRequest() (*http.Request, error) {
  return http.NewRequest("GET", "/api/crud/site/", nil)
}
func siteListRequestLimited() (*http.Request, error) {
  return http.NewRequest("GET", "/api/crud/site/?limit=1&offset=1", nil)
}
func siteListRequestLimited2() (*http.Request, error) {
  return http.NewRequest("GET", "/api/crud/site/?limit=2&offset=2", nil)
}

func TestList(t *testing.T) {
  if err := startToGolden("site-list", siteListRequest);
       err != nil {
    t.Error(err.Error())
  }
}

func TestListLimit(t *testing.T) {
  repos, handler, err := startToSetup()
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  if err := setupToGolden(repos, handler, "site-list-limit-1", siteListRequestLimited); err != nil {
    t.Fatal(err.Error())
  }
  if err := setupToGolden(repos, handler, "site-list-limit-2", siteListRequestLimited2); err != nil {
    t.Fatal(err.Error())
  }
}
