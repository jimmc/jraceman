// The conn package provides a database connector that can be either a
// direct connection to a database, or a transaction.
package conn

import (
  "database/sql"
)

// DB defines the methods we use to access our database that are common
// to sql.DB and sql.Tx, so that we can use either one in a Repos struct
// and in SQL statements.
type DB interface {
  Exec(query string, args ...any) (sql.Result, error)
  Query(query string, args ...any) (*sql.Rows, error)
  QueryRow(query string, args ...any) *sql.Row
}
