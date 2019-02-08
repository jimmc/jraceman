package dbrepo

import (
  "database/sql"
  "fmt"
  "strings"

  // TODO - We only need structsql for requireOneResult, maybe that can go elsewhere
  "github.com/jimmc/jracemango/dbrepo/structsql"
)

// dbRowRepo implements the RowRepo interface for use by Importer.
type dbRowRepo struct {
  db *sql.DB    // One or the other of db and tx must be filled.
  tx *sql.Tx
}

func NewRowRepo(dbrepos *Repos) *dbRowRepo {
  return &dbRowRepo{
    db: dbrepos.db,
  }
}

func NewRowRepoWithTx(dbrepos *Repos) (*dbRowRepo, error) {
  tx, err := dbrepos.db.Begin()
  if err != nil {
    return nil, err
  }
  return &dbRowRepo{
    tx: tx,
  }, nil
}

func (r *dbRowRepo) Commit() error {
  if r.tx == nil {
    return fmt.Errorf("No transaction to commit")
  }
  return r.tx.Commit()
}

func (r *dbRowRepo) Rollback() error {
  if r.tx == nil {
    return fmt.Errorf("No transaction to rollback")
  }
  return r.tx.Rollback()
}

func (r *dbRowRepo) Read(table string, columns []string, ID string) ([]interface{}, error) {
  values := make([]interface{}, len(columns))
  targets := make([]interface{}, len(columns))
  for i := 0; i < len(values); i++ {
    targets[i] = &values[i]
  }
  selSql := "SELECT " + strings.Join(columns, ",") + " from " + table + " where id=?;"
  err := r.queryRow(selSql, ID).Scan(targets...)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil           // No data and no error
    } else {
      return nil, fmt.Errorf("error retrieving existing row %s[%s]: %v; sql=%s",
          table, ID, err, selSql)
    }
  }
  // Convert to the expected common data types
  for i, v := range values {
    switch vv := v.(type) {
    case []uint8:
      values[i] = string(vv)
    case int64:
      values[i] = int(vv)
    }
  }
  return values, nil
}

func (r *dbRowRepo) Insert(table string, columns[]string, values []interface{}, ID string) error {
  insSql := "INSERT into " + table + "(" + strings.Join(columns, ",") + ") values(" +
      strings.Repeat("?,", len(columns) - 1) + "?);"
  res, err := r.exec(insSql, values...)
  return structsql.RequireOneResult(res, err, "Inserted", table, ID)
}

func (r *dbRowRepo) Update(table string, columns[]string, values []interface{}, ID string) error {
  insSql := "UPDATE " + table + " set " + strings.Join(columns, " = ?, ") + " = ? where id = ?;"
  values = append(values, ID)
  res, err := r.exec(insSql, values...)
  return structsql.RequireOneResult(res, err, "Updated", table, ID)
}

func (r *dbRowRepo) queryRow(query string, args ...interface{}) *sql.Row {
  if r.tx != nil {
    return r.tx.QueryRow(query, args...)
  }
  return r.db.QueryRow(query, args...)
}

func (r *dbRowRepo) exec(sqlstr string, values ...interface{}) (sql.Result, error) {
  if r.tx != nil {
    return r.tx.Exec(sqlstr, values...)
  }
  return r.db.Exec(sqlstr, values...)
}
