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
// For a single test, the typical calling sequence is:
//   r := NewTester(handlerCreateFunc)
//   r.Run(t, basename, callback)
// For multiple tests, maintaining the database state across tests as it changes:
//   r := NewTester(handlerCreateFunc)
//   r.Init()
//   r.RunTestWith(t, basename, callback)
//   r.RunTestWith(t, basename2, callback2)
//   r.Close()
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

// NewTester creates a new instance of a Tester that will use the specified
// function to create an http.Handler.
func NewTester(createHandler func(r *Tester) http.Handler) *Tester {
  r := &Tester{}
  r.CreateHandler = createHandler
  return r
}

// SetBaseNameAndCallback resets the basename and callback of the Tester in preparation for running a test.
func (r *Tester) SetBaseNameAndCallback(basename string, callback func() (*http.Request, error)) {
  r.BaseName = basename
  r.Callback = callback
}

// SetupFilePath returns the complete path to the setup file.
func (r *Tester) SetupFilePath() string {
  return r.GetFilePath(r.SetupPath, r.SetupBaseName, "setup")
}

// Init does all of the setup from base.Tester, and sets up the database.
func (r *Tester) Init(t *testing.T) {
  t.Helper()
  if err := r.Tester.Init(); err != nil {
    t.Fatalf("Error in base.Tester.Init: %v", err)
  }
  repos, err := dbtest.ReposEmpty()
  if err != nil {
    t.Fatalf("Error creating Repos: %v", err)
  }
  r.repos = repos
}

// Arrange loads a setup file into the database. The caller should update the
// SetupBaseName or SetupPath fields before calling this function.
func (r *Tester) Arrange() error {
  setupfilepath := r.SetupFilePath()
  if err := goldendb.LoadSetupFile(r.repos.DB(), setupfilepath); err != nil {
    r.repos.Close()
    return fmt.Errorf("error loading setup file %v: %v", setupfilepath, err)
  }
  if err := r.Tester.Arrange(); err != nil {
    return err
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
  return r.Tester.Assert()
}

// RunTest loads a new setup file, runs the test, and compares the output.
// This function can be used in a test where there are multiple files to be loaded.
// A typical calling sequence for that scenario is to call Init,
// then RunTest multiple times, and finally Close when done with all tests.
func (r *Tester) RunTest(t *testing.T) {
  t.Helper()
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error loading setup file: %v", err)
  }
  if err := r.Act(); err != nil {
    t.Fatalf("Error running test: %v", err)
  }
  if err := r.Assert(); err != nil {
    t.Fatal(err)
  }
}

// Close closes our database.
func (r *Tester) Close() {
  if r.repos != nil {
    r.repos.Close()
  }
}

// RunTestWith runs a test using the specified basename and callback.
// This can be used multiple times within a Tester. The database state is maintained across tests,
// allowing a sequence of calls that builds up and modifies a database.
func (r *Tester) RunTestWith(t *testing.T, basename string, callback func() (*http.Request, error)) {
  t.Helper()
  r.SetBaseNameAndCallback(basename, callback)
  r.RunTest(t)
}

// Run initializes the tester, runs a test, and closes it, calling Fatalf on any error.
func (r *Tester) Run(t *testing.T, basename string, callback func() (*http.Request, error)) {
  t.Helper()
  r.Init(t)
  r.RunTestWith(t, basename, callback)
  r.Close()
}
