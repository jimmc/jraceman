package structsql

import (
  "reflect"
  "strconv"
  "strings"

  "github.com/golang/glog"
)

// Placeholder is ? for MySQL,$N for PostgreSQL,
// SQLite uses either of those, Oracle is :param1

// SelectSql generates an SQL SELECT statement
// using the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
// There is no semicolon at the end of the sql, so the caller can append more
// sql to it such as a where clause.
func SelectSql(tableName string, entity interface{}) (string, []interface{}) {
  val := reflect.Indirect(reflect.ValueOf(entity))
  typ := val.Type()
  numFields := typ.NumField()
  columnNames := make([]string, numFields)
  targets := make([]interface{}, numFields)
  for i := 0; i < numFields; i++ {
    field := typ.Field(i)
    columnName := strings.ToLower(field.Name)
    columnNames[i] = columnName
    targets[i] = val.Field(i).Addr().Interface()
  }
  sql := "SELECT " + strings.Join(columnNames, ",") + " from " + tableName
  return sql, targets
}

// FindByIDSql generates an SQL QUERY statement,
// with a WHERE clause limiting it to a matching id,
// using the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
func FindByIDSql(tableName string, entity interface{}) (string, []interface{}) {
  sql, targets := SelectSql(tableName, entity)
  sql = sql + " where id=?;"
  glog.V(1).Infof("FindByIDSql: %v\n", sql)
  return sql, targets
}

// ListSql generates an SQL QUERY statement,
// with OFFSET and LIMIT clauses,
// using the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
func ListSql(tableName string, entity interface{}, offset, limit int) (string, []interface{}) {
  sql, targets := SelectSql(tableName, entity)
  if limit != 0 {
    sql = sql + " limit " + strconv.Itoa(limit)
  }
  if offset != 0 {
    sql = sql + " offset " + strconv.Itoa(offset)
  }
  glog.V(1).Infof("ListSql: %v\n", sql)
  return sql, targets
}
