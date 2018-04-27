package dbrepo

import (
  "database/sql"
  "fmt"
  "log"
  "reflect"
  "strings"
)

// StdCreateTableSqlFromStruct generates an SQL CREATE TABLE command using
// the fields of the given struct. All field names are converted to lower case.
func stdCreateTableSqlFromStruct(tableName string, entity interface{}) string {
  val := reflect.Indirect(reflect.ValueOf(entity))
  typ := val.Type()
  numFields := typ.NumField()
  columnSpecs := make([]string, numFields)
  for i := 0; i < numFields; i++ {
    field := typ.Field(i)
    columnName := strings.ToLower(field.Name)
    goTypeName := field.Type.String()     // string, *string, int
    goTypeName = strings.TrimPrefix(goTypeName, "*")
    columnType := goTypeName            // TODO - convert as required
    columnSpec := columnName + " " + columnType
    if columnName == "id" {
      columnSpec = columnSpec + " primary key"
    }
    columnSpecs[i] = columnSpec
  }
  sql := "CREATE TABLE " + tableName + "(" + strings.Join(columnSpecs, ",") + ");"
  return sql
}

// StdFindByIDSqlFromStruct generatesand SQL QUERY statement using
// thefields of the give struct.
func stdFindByIDSqlFromStruct(tableName string, entity interface{}) (string, []interface{}) {
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
  sql := "SELECT " + strings.Join(columnNames, ",") + " from " + tableName + " where id=?"
  return sql, targets
}

// RequireOneResult gets the result of sql.Stmt.Exec and verifies that it
// affected exactly one row, which should be the case for operations that
// use the entity ID. The action string is used in error messages, and
// should be capitalized and past tense, such as "Deleted".
// The entityType should the name of the entity, such as "site".
func requireOneResult(res sql.Result, err error, action, entityType, ID string) error {
  if err != nil {
    return err
  }
  rowCnt, err := res.RowsAffected()
  if err != nil {
    return err
  }
  if rowCnt == 0 {
    return fmt.Errorf("%s no %s rows for ID %s", action, entityType, ID)
  }
  if rowCnt > 1 {
    return fmt.Errorf("%s %d %s rows for ID %s", action, rowCnt, entityType, ID)
  }
  return nil
}

// ColumnsUpdateStringAndValues generates a string for the column-and-values portion
// of an SQL update statement, in the form "col1 = ?, col2 = ?", and also returns
// an array of values that correspond to those columns.
func columnsUpdateStringAndVals(mods map[string]interface{}) (string, []interface{}) {
  var keys []string
  var vals []interface{}
  for k, v := range mods {
    if strings.HasPrefix(k, ".") {
      k = strings.TrimPrefix(k, ".")
    }
    // By default, we use all lowercase names for database columns.
    k = strings.ToLower(k)
    keys = append(keys, k + " = ?")
    vals = append(vals, v)
  }
  return strings.Join(keys, ", "), vals
}

// ModsToSql takes a map of the modifications from Diffs and generates the
// sql string and values to be executed to perform the update.
func modsToSql(table string, mods map[string]interface{}, ID string) (string, []interface{}) {
  log.Printf("mods = %v", mods)
  kvString, vals := columnsUpdateStringAndVals(mods)
  updateSql := "update " + table + " set " + kvString + " where id = ?"
  vals = append(vals, ID)
  log.Printf("updateSql = %q, vals = %v", updateSql, vals)
  return updateSql, vals
}
