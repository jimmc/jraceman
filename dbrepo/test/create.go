// Package test contains test helper functions dealing with a database.
package test

import (
  "database/sql"
  "fmt"
  "io/ioutil"
  "os"

  "github.com/jimmc/jraceman/dbrepo"

  "github.com/golang/glog"

  goldendb "github.com/jimmc/golden/db"

  _ "github.com/mattn/go-sqlite3"
)

// Struct TestRecord matches the format of the test table
// created by DbWithTestTable.
type TestRecord struct {
  ID string;
  N int;
  S string;
  S2 *string;
}

// DbWIthTestTable creates a test database with one small table
// called test.
func DbWithTestTable() (*sql.DB, error) {
  return goldendb.DbWithSetupString(`
CREATE table test(id string, n int, s string, s2 string);

INSERT into test(id, n, s, s2)
values('T1', 1, 'a', 'A'), ('T2', 2, 'b', null), ('T3', 3, 'c', 'C');
`)
}

// ReposEmpty creates an empty test Repos.
// The second return value is the cleanup function. Tests should call
// this function to ensure that everything is properly cleaned up
// at the end of each test.
func ReposEmpty() (*dbrepo.Repos, func(), error) {
  // The in-memory database has some quirks that make it harder to use for testing.
  // See the comments in ReposEmptyInMemory. Using a TempFile database is still pretty
  // fast for our tests, and doesn't have those problems.
  // return ReposEmptyInMemory()
  return ReposEmptyTempFile()
}

// ReposEmptyInMemory create an empty in-memory Repos for
// use with unit tests.
func ReposEmptyInMemory() (*dbrepo.Repos, func(), error) {
  // If we specify ":memory": to sqlite3, each connection opens a new
  // in-memory database. This appears to happen when you try to issue
  // a query on a connection while another query is still open, which
  // makes tests fail with a "no such table" error. According to the FAQ
  // (https://github.com/mattn/go-sqlite3/blob/master/README.md#faq)
  // a workaround is to file "file::memory:?cache=shared", which will
  // make all connections point to the same database.
  // Unfortunately, this also means the next time an in-memory database
  // is opened, it opens the same database, so the cleanup is not done
  // properly. A potential workaround for this problem is given here:
  // https://stackoverflow.com/questions/21547616/how-do-i-completely-clear-a-sqlite3-database-without-deleting-the-database-file
  nop := func(){}
  repos, err :=  dbrepo.Open("sqlite3:file::memory:?cache=shared")
  if err != nil {
    return nil, nop, fmt.Errorf("failed to open in-memory database: %w", err)
  }
  cleanup := func() {
    repos.Close()
  }
  return repos, cleanup, err
}

// ReposEmptyTempFile create an empty Repos database using a temp
// file for use with unit tests.
func ReposEmptyTempFile() (*dbrepo.Repos, func(), error) {
  nop := func(){}
  tmpf, err := ioutil.TempFile("", "jrdbtest")
  if err != nil {
    return nil, nop, fmt.Errorf("failed to open tmp file for database: %w", err)
  }
  tmpf.Close()
  dbr, err := dbrepo.Open("sqlite3:file:"+tmpf.Name())
  if err != nil {
    os.Remove(tmpf.Name())
    return nil, nop, fmt.Errorf("failed to open database on tmp file %q: %w",
        tmpf.Name(), err)
  }
  fn := tmpf.Name()
  cleanup := func() {
    dbr.Close()
    os.Remove(fn)
  }
  return dbr, cleanup, nil
}

// ReposAndLoadFile creates a test Repos with the default
// set of tables, then imports the specified JRaceman-format file.
func ReposAndLoadFile(setupfile string) (*dbrepo.Repos, func(), error) {
  repos, cleanup, err := ReposEmpty()
  if err != nil {
    cleanup()
    return nil, nil, fmt.Errorf("failed to open repository: %v", err)
  }

  err = repos.CreateTables()
  if err != nil {
    cleanup()
    return nil, nil, fmt.Errorf("failed to create repository tables: %v", err)
  }

  glog.Infof("Importing from %s\n", setupfile)
  counts, err := repos.ImportFile(setupfile)
  if err != nil {
    cleanup()
    return nil, nil, fmt.Errorf("error importing from %s: %v", setupfile, err)
  }
  glog.Infof("Import done: inserted %d, updated %d, unchanged %d records\n",
      counts.Inserted(), counts.Updated(), counts.Unchanged())
  return repos, cleanup, nil
}

// ReposAndSqlFile creates a test Repos and loads
// the specified SQL file. It does not create the standard tables.
func EmptyReposAndSqlFile(setupfile string) (*dbrepo.Repos, func(),  error) {
  repos, cleanup, err := ReposEmpty()
  if err != nil {
    cleanup()
    return nil, nil, fmt.Errorf("failed to open repository: %v", err)
  }

  glog.Infof("Loading SQL from %s\n", setupfile)
  db, ok := repos.DB().(*sql.DB)
  if !ok {
    cleanup()
    return nil, nil, fmt.Errorf("repos.DB() is not *sql.DB")
  }
  err = goldendb.LoadSetupFile(db, setupfile)
  if err != nil {
    cleanup()
    return nil, nil, fmt.Errorf("error importing from %s: %v", setupfile, err)
  }
  return repos, cleanup, nil
}
