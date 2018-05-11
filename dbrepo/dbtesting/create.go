// Package dbtesting contains test helper functions dealing with a database.
package dbtesting

import (
  "database/sql"
  "io/ioutil"

  "github.com/jimmc/jracemango/dbrepo"
  "github.com/jimmc/jracemango/dbrepo/strsql"

  _ "github.com/mattn/go-sqlite3"
)

func EmptyDb() (*sql.DB, error) {
  return sql.Open("sqlite3", ":memory:")
}

func DbWithSetupFile(filename string) (*sql.DB, error) {
  setupString, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return DbWithSetupString(string(setupString))
}

func DbWithSetupString(setupSql string) (*sql.DB, error) {
  db, err := EmptyDb()
  if err != nil {
    return nil, err
  }
  err = strsql.ExecMulti(db, setupSql)
  if err != nil {
    db.Close()
    return nil, err
  }
  return db, nil
}

// Struct TestRecord matches the format of the test table
// created by DbWithTestTable.
type TestRecord struct {
  ID string;
  N int;
  S string;
  S2 *string;
}

func DbWithTestTable() (*sql.DB, error) {
  return DbWithSetupString(`
CREATE table test(id string, n int, s string, s2 string);

INSERT into test(id, n, s, s2)
values('T1', 1, 'a', 'A'), ('T2', 2, 'b', null), ('T3', 3, 'c', 'C');
`)
}

func ReposEmpty() (*dbrepo.Repos, error) {
  return dbrepo.Open("sqlite3::memory:")
}
