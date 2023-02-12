package compat

import (
  "database/sql"
)

// DBorTx defines the methods we use to access our database that are common
// to sql.DB and sql.Tx, so that we can use either one in a Repos struct
// and in SQL statements.
type DBorTx interface {
  Exec(query string, args ...any) (sql.Result, error)
  Query(query string, args ...any) (*sql.Rows, error)
  QueryRow(query string, args ...any) *sql.Row
}
