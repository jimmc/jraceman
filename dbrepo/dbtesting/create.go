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

func DbWithTestTable() (*sql.DB, error) {
  return DbWithSetupString(`
CREATE table test(id string, n int, s string);

INSERT into test(id, n, s) values('T1', 1, 'a'), ('T2', 2, 'b'), ('T3', 3, 'c');
`)
}

func ReposEmpty() (*dbrepo.Repos, error) {
  return dbrepo.Open("sqlite3::memory:")
}
