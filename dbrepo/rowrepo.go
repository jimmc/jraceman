package dbrepo

import (
  "database/sql"
  "fmt"
  "strings"
)

// dbRowRepo implements the RowRepo interface for use by Importer.
type dbRowRepo struct {
  db *sql.DB
}

func (r *dbRowRepo) Read(table string, columns []string, ID string) ([]interface{}, error) {
  targets := make([]interface{}, len(columns))
  selSql := "SELECT " + strings.Join(columns, ",") + " from " + table + " where id=?;"
  err := r.db.QueryRow(selSql, ID).Scan(targets...)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil           // No data and no error
    } else {
      return nil, fmt.Errorf("error retrieving existing row %s[%s]: %v",
          table, ID, err)
    }
  }
  return targets, nil
}

func (r *dbRowRepo) Insert(table string, columns[]string, values []interface{}, ID string) error {
  insSql := "INSERT into " + table + "(" + strings.Join(columns, ",") + ") values(" +
      strings.Repeat("?,", len(columns) - 1) + "?);"
  res, err := r.db.Exec(insSql, values...)
  return requireOneResult(res, err, "Inserted", table, ID)
}

func (r *dbRowRepo) Update(table string, columns[]string, values []interface{}, ID string) error {
  insSql := "UPDATE " + table + " set " + strings.Join(columns, " = ?, ") + " = ? where id = ?;"
  values = append(values, ID)
  res, err := r.db.Exec(insSql, values...)
  return requireOneResult(res, err, "Updated", table, ID)
}
