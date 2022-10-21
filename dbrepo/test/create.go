// Package test contains test helper functions dealing with a database.
package test

import (
  "database/sql"

  "github.com/jimmc/jraceman/dbrepo"

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

func DbWithTestTable() (*sql.DB, error) {
  return goldendb.DbWithSetupString(`
CREATE table test(id string, n int, s string, s2 string);

INSERT into test(id, n, s, s2)
values('T1', 1, 'a', 'A'), ('T2', 2, 'b', null), ('T3', 3, 'c', 'C');
`)
}

func ReposEmpty() (*dbrepo.Repos, error) {
  return dbrepo.Open("sqlite3::memory:")
}
