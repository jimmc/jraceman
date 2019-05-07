package test

import (
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "os"
  "testing"

  "github.com/jimmc/jracemango/dbrepo"
  dbtest "github.com/jimmc/jracemango/dbrepo/test"

  goldenbase "github.com/jimmc/golden/base"
  goldendb "github.com/jimmc/golden/db"
)

// Tester provides the structure for running API unit tests.
type Tester struct {
  goldenbase.Tester

  // Base name for the test setup file; if not set, uses BaseName.
  SetupBaseName string
  // Path to the test setup file; if not set, uses SetupBaseName.
  SetupPath string

  CreateHandler func(r *Tester) http.Handler
  Callback func() (*http.Request, error)

  repos *dbrepo.Repos
}

// NewTester creates a new instance of a Tester that will call the specified
// callback as the test function.
func NewTester(basename string, createHandler func(r *Tester) http.Handler, callback func() (*http.Request, error)) *Tester {
  r := &Tester{}
  r.BaseName = basename
  r.CreateHandler = createHandler
  r.Callback = callback
  return r
}

// SetupFilePath returns the complete path to the setup file.
func (r *Tester) SetupFilePath() string {
  return r.GetFilePath(r.SetupPath, r.SetupBaseName, "setup")
}

// Setup does all of the setup from base.Tester, and sets up the db
// and loads the setup file as calculated from the Tester.
func (r *Tester) Setup() error {
  if err := r.SetupDb(); err != nil {
    return err
  }
  return r.LoadSetupFile()
}

// Setup does all of the setup from base.Tester, and sets up the database.
func (r *Tester) SetupDb() error {
  if err := r.Tester.Setup(); err != nil {
    return err
  }
  repos, err := dbtest.ReposEmpty()
  if err != nil {
    return err
  }
  r.repos = repos
  return nil
}

// LoadSetupFile loads a setup file into the database. The caller should update the
// SetupBaseName or SetupPath fields before calling this function.
func (r *Tester) LoadSetupFile() error {
  setupfilepath := r.SetupFilePath()
  if err := goldendb.LoadSetupFile(r.repos.DB(), setupfilepath); err != nil {
    r.repos.Close()
    return fmt.Errorf("error loading setup file %v: %v", setupfilepath, err)
  }
  return nil
}

// Act sets up the handler, calls the request, and records the result to the output file.
func (r *Tester) Act() error {
  handler := r.CreateHandler(r)

  req, err := r.Callback()
  if err != nil {
    return fmt.Errorf("error calling callback in Tester.Act: %v", err)
  }

  rr := httptest.NewRecorder()
  handler.ServeHTTP(rr, req)

  if got, want := rr.Code, http.StatusOK; got != want {
    return fmt.Errorf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }

  body := rr.Body.Bytes()
  if len(body) == 0 {
    return errors.New("response body should not be empty")
  }

  outfilepath := r.OutFilePath()
  os.Remove(outfilepath)
  if err := ioutil.WriteFile(outfilepath, body, 0644); err != nil {
    return err
  }
  return nil
}

// Assert closes the output file and compares it to the golden file.
func (r *Tester) Assert() error {
  return r.Tester.Finish()
}

// LoadActAssert loads a new setup file, runs the test, and compares the output.
// This function can be used in a test where there are multiple files to be loaded.
// A typical calling sequence for that scenario is to call SetupDb,
// then LoadActAssert multiple times, and finally Close() when done with all tests.
func (r *Tester) LoadActAssert() error {
  if err := r.LoadSetupFile(); err != nil {
    return err
  }
  if err := r.Act(); err != nil {
    return err
  }
  if err := r.Assert(); err != nil {
    return err
  }
  return nil
}

// LoadActAssertT runs LoadActAssert and calls Fatalf on error.
func (r *Tester) LoadActAssertT(t *testing.T) {
  if err := r.LoadActAssert(); err != nil {
    t.Fatalf("Error running LoadActAssert: %v", err)
  }
}

func (r *Tester) Close() {
  if r.repos != nil {
    r.repos.Close()
  }
}

// Finish closes the database and the output file and checks the output.
func (r *Tester) Finish() error {
  r.Close()
  return r.Assert()
}
