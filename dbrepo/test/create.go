// Package test contains test helper functions dealing with a database.
package test

import (
  "database/sql"
  "fmt"

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

// DbWIthTestTable creates an in-memory database with one small table
// called test.
func DbWithTestTable() (*sql.DB, error) {
  return goldendb.DbWithSetupString(`
CREATE table test(id string, n int, s string, s2 string);

INSERT into test(id, n, s, s2)
values('T1', 1, 'a', 'A'), ('T2', 2, 'b', null), ('T3', 3, 'c', 'C');
`)
}

// ReposEmpty creates an in-memory and empty Repos.
func ReposEmpty() (*dbrepo.Repos, error) {
  return dbrepo.Open("sqlite3::memory:")
}

// ReposAndLoadFile creates an in-memory Repos with the default
// set of tables, then imports the specified JRaceman-format file.
func ReposAndLoadFile(setupfile string) (*dbrepo.Repos, error) {
  repos, err := ReposEmpty()
  if err != nil {
    return nil, fmt.Errorf("failed to open repository: %v", err)
  }

  err = repos.CreateTables()
  if err != nil {
    return nil, fmt.Errorf("failed to create repository tables: %v", err)
  }

  glog.Infof("Importing from %s\n", setupfile)
  counts, err := repos.ImportFile(setupfile)
  if err != nil {
    return nil, fmt.Errorf("error importing from %s: %v", setupfile, err)
  }
  glog.Infof("Import done: inserted %d, updated %d, unchanged %d records\n",
      counts.Inserted(), counts.Updated(), counts.Unchanged())
  return repos, nil
}

// ReposAndSqlFile creates an in-memory Repos and loads
// the specified SQL file. It does not create the standard tables.
func EmptyReposAndSqlFile(setupfile string) (*dbrepo.Repos, error) {
  repos, err := ReposEmpty()
  if err != nil {
    return nil, fmt.Errorf("failed to open repository: %v", err)
  }

  glog.Infof("Loading SQL from %s\n", setupfile)
  err = goldendb.LoadSetupFile(repos.DB(), setupfile)
  if err != nil {
    return nil, fmt.Errorf("error importing from %s: %v", setupfile, err)
  }
  return repos, nil
}
