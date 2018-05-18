package test

import (
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "os"

  "github.com/jimmc/jracemango/api/crud"
  "github.com/jimmc/jracemango/dbrepo"
  dbtest "github.com/jimmc/jracemango/dbrepo/test"
)

func StartToGolden(basename string, callback func() (*http.Request, error)) error {
  repos, handler, err := StartToSetup()
  if err != nil{
    return err
  }
  defer repos.Close()

  return SetupToGolden(repos, handler, basename, callback)
}

// startToSetup initializes the database and the http handler.
func StartToSetup() (*dbrepo.Repos, http.Handler, error) {
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

func SetupToGolden(repos *dbrepo.Repos, handler http.Handler, basename string,
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
    return fmt.Errorf("response status: got %d, want %d\nBody: %v", got, want, rr.Body.String())
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
