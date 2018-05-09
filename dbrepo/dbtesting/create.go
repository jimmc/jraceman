// Package dbtesting contains test helper functions dealing with a database.
package dbtesting

import (
  "database/sql"

  "github.com/jimmc/jracemango/dbrepo"

  _ "github.com/mattn/go-sqlite3"
)

func EmptyDb() (*sql.DB, error) {
  return sql.Open("sqlite3", ":memory:")
}

func ReposEmpty() (*dbrepo.Repos, error) {
  return dbrepo.Open("sqlite3::memory:")
}

func CreateAndPopulateTestTable(db *sql.DB) error {
  createSql := "CREATE table test(n int, s string);"
  if _, err := db.Exec(createSql); err != nil {
    return err
  }

  insertSql := "INSERT into test(n, s) values(1, 'a'), (2, 'b'), (3, 'c');"
  if _, err := db.Exec(insertSql); err != nil {
    return err
  }

  return nil
}
