package test

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jraceman/dbrepo"

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
