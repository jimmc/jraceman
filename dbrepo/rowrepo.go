package dbrepo

import (
  "database/sql"
  "fmt"
  "strings"

  "github.com/jimmc/jraceman/dbrepo/conn"
  // TODO - We only need structsql for requireOneResult, maybe that can go elsewhere
  "github.com/jimmc/jraceman/dbrepo/structsql"
)

// dbRowRepo implements the RowRepo interface for use by Importer.
type dbRowRepo struct {
  db conn.DB    // One or the other of db and tx must be filled.
  tx *sql.Tx
}

func NewRowRepo(dbrepos *Repos) *dbRowRepo {
  return &dbRowRepo{
    db: dbrepos.db,
  }
}

func NewRowRepoWithTx(dbrepos *Repos) (*dbRowRepo, error) {
  db, ok := dbrepos.db.(*sql.DB)
  if !ok {
    return nil, fmt.Errorf("Can't start a transaction on a non-DB")
  }
  tx, err := db.Begin()
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
  return r.ReadByKey(table, columns, "id", ID)
}

func (r *dbRowRepo) ReadByKey(table string, columns []string, keyName, key string) ([]interface{}, error) {
  values := make([]interface{}, len(columns))
  targets := make([]interface{}, len(columns))
  for i := 0; i < len(values); i++ {
    targets[i] = &values[i]
  }
  selSql := "SELECT " + strings.Join(columns, ",") + " from " + table + " where " + keyName + "=?;"
  err := r.queryRow(selSql, key).Scan(targets...)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil           // No data and no error
    } else {
      return nil, fmt.Errorf("error retrieving existing row %s[%s]: %v; sql=%s",
          table, key, err, selSql)
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

func (r *dbRowRepo) Insert(table string, columns[]string, values []interface{}, key string) error {
  insSql := "INSERT into " + table + "(" + strings.Join(columns, ",") + ") values(" +
      strings.Repeat("?,", len(columns) - 1) + "?);"
  res, err := r.exec(insSql, values...)
  return structsql.RequireOneResult(res, err, "Inserted", table, key)
}

func (r *dbRowRepo) Update(table string, columns[]string, values []interface{}, ID string) error {
  return r.UpdateByKey(table, columns, values, "id", ID)
}

func (r *dbRowRepo) UpdateByKey(table string, columns[]string, values []interface{}, keyName, key string) error {
  insSql := "UPDATE " + table + " set " + strings.Join(columns, " = ?, ") + " = ? where " + keyName + " = ?;"
  values = append(values, key)
  res, err := r.exec(insSql, values...)
  return structsql.RequireOneResult(res, err, "Updated", table, key)
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
