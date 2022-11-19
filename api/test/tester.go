package test

import (
  "context"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "os"
  "testing"

  "github.com/jimmc/jraceman/dbrepo"
  dbrepotest "github.com/jimmc/jraceman/dbrepo/test"

  authusers "github.com/jimmc/auth/users"
  authperms "github.com/jimmc/auth/permissions"
  goldenbase "github.com/jimmc/golden/base"
  goldenhttpdb "github.com/jimmc/golden/httpdb"
)

// Tester provides the structure for running API unit tests.
// For a single test, the typical calling sequence is:
//   r := NewTester(handlerCreateFunc)
//   r.Run(basename, callback)
// For multiple tests, maintaining the database state across tests as it changes:
//   r := NewTester(handlerCreateFunc)
//   r.Init()
//   r.RunTestWith(basename, callback)
//   r.RunTestWith(basename2, callback2)
//   r.Close()
type Tester struct {
  goldenhttpdb.Tester

  repos *dbrepo.Repos
}

// NewTester creates a new instance of a Tester that will use the specified
// function to create an http.Handler.
func NewTester(createHandler func(r *Tester) http.Handler) *Tester {
  r := &Tester{}
  r.Tester.CreateHandler = func(*goldenhttpdb.Tester) http.Handler {
    return createHandler(r)
  }
  return r
}

// Init does all of the setup from goldenhttpdb.Tester, and sets up the database.
func (r *Tester) Init() error {
  if err := r.Tester.Init(); err != nil {
    return fmt.Errorf("error in base.Tester.Init: %v", err)
  }
  repos, err := dbrepo.OpenDB(r.Tester.DB)
  if err != nil {
    return fmt.Errorf("error creating Repos: %v", err)
  }
  r.repos = repos
  return nil
}

// Close closes our database.
func (r *Tester) Close() error {
  if r.repos != nil {
    r.repos.Close()
  }
  r.Tester.Close()
  return nil
}

// Run initializes the tester, runs a test, and closes it, calling Fatalf on any error.
func RunOneWith(r *Tester, basename string, callback func() (*http.Request, error)) error {
  if err := r.Init(); err != nil {
    return err
  }
  if err := goldenhttpdb.RunTestWith(r, basename, callback); err != nil {
    return err
  }
  return r.Close()
}

// AddTestUser adds a test user to the request, with the specified permissions.
func AddTestUser(r *http.Request, permstr string) *http.Request {
  username := "testuser"
  saltword := "not-used"
  perms := authperms.FromString(permstr) // permstr is a space-separated list of permission.
  user := authusers.NewUser(username, saltword, perms)
  cwv := context.WithValue(r.Context(), "AuthUser", user)
  return r.WithContext(cwv)
}

// Returns the database repo and a cleanup function to close the repo.
func RequireDatabaseWithSqlFile(t *testing.T, filebase string) (*dbrepo.Repos, func()) {
  filename := "testdata/" + filebase + ".setup"
  repos, err := dbrepotest.EmptyReposAndSqlFile(filename)
  if err != nil {
    t.Fatalf("Error creating database: %v", err)
  }
  cleanup := func() {
    repos.Close()
  }
  return repos, cleanup
}

func RequireHttpSuccess(t *testing.T, req *http.Request, rr *httptest.ResponseRecorder) {
  t.Helper()
  if got, want := rr.Code, http.StatusOK; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func RequireBodyMatchesGolden(t *testing.T, rr *httptest.ResponseRecorder, filebase string) {
  t.Helper()
  body := rr.Body.Bytes()
  if len(body) == 0 {
    t.Fatal("Response body should not be empty")
  }

  outfilepath := "testdata/" + filebase + ".out"
  goldenfilepath := "testdata/" + filebase + ".golden"
  os.Remove(outfilepath)
  if err := ioutil.WriteFile(outfilepath, body, 0644); err != nil {
    t.Fatalf("Error writing body to file %q: %v", outfilepath, err)
  }
  err := goldenbase.CompareOutToGolden(outfilepath, goldenfilepath)
  if err != nil {
    t.Fatalf("Output did not match: %v", err)
  }
}
