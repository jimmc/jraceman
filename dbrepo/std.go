package dbrepo

import (
  "database/sql"
  "fmt"
  "log"
  "reflect"
  "strings"
)

// Placeholder is ? for MySQL,$N for PostgreSQL,
// SQLite uses either of those, Oracle is :param1

// StdCreateTableSqlFromStruct generates an SQL CREATE TABLE command using
// the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case for the column name.
//   * int and string fields are declared as that same type column.
//   * The id field is declared as primary key.
//   * Non-pointer fields are declared as not null.
func stdCreateTableSqlFromStruct(tableName string, entity interface{}) string {
  val := reflect.Indirect(reflect.ValueOf(entity))
  typ := val.Type()
  numFields := typ.NumField()
  columnSpecs := make([]string, numFields)
  for i := 0; i < numFields; i++ {
    field := typ.Field(i)
    columnName := strings.ToLower(field.Name)
    goTypeName := field.Type.String()     // string, *string, int
    isPointer := strings.HasPrefix(goTypeName, "*")
    goTypeName = strings.TrimPrefix(goTypeName, "*")
    columnType := goTypeName            // TODO - convert as required
    columnSpec := columnName + " " + columnType
    if columnName == "id" {
      columnSpec = columnSpec + " primary key"
    } else if !isPointer {
      columnSpec = columnSpec + " not null"
    }
    columnSpecs[i] = columnSpec
  }
  sql := "CREATE TABLE " + tableName + "(" + strings.Join(columnSpecs, ", ") + ");"
  log.Printf("stdCreateTableSql: %v\n", sql)
  return sql
}

// StdFindByIDSqlFromStruct generates an SQL QUERY statement using
// the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
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
  sql := "SELECT " + strings.Join(columnNames, ",") + " from " + tableName + " where id=?;"
  log.Printf("stdFindByIDSql: %v\n", sql)
  return sql, targets
}

// StdInsertSqlFromStruct generates an SQL INSERT statement using
// the fields of the given struct. For each field in the struct:
//   * If the field is a nil pointer, it is ignore.
//   * The field name is converted to lower case.
func stdInsertSqlFromStruct(tableName string, entity interface{}) (string, []interface{}) {
  val := reflect.Indirect(reflect.ValueOf(entity))
  typ := val.Type()
  numFields := typ.NumField()
  columnNames := make([]string, 0)
  values := make([]interface{}, 0)
  placeHolders := make([]string, 0)
  for i := 0; i < numFields; i++ {
    fieldType := typ.Field(i)
    columnName := strings.ToLower(fieldType.Name)
    vf := val.Field(i)
    if !(vf.Kind() == reflect.Ptr && vf.IsNil()) {  // Omit nil pointers
      if vf.Kind() == reflect.Ptr {
        vf = vf.Elem()  // Dereference the pointer to get the value
      }
      fieldValue := vf.Interface()
      columnNames = append(columnNames, columnName)
      values = append(values, fieldValue)
      placeHolders = append(placeHolders, "?")
    }
  }
  sql := "INSERT into " + tableName + "(" + strings.Join(columnNames, ",") + ")" +
      " values (" + strings.Join(placeHolders,",") + ");"
  log.Printf("stdInsertSql: %v\n", sql)
  log.Printf("  values: %v\n", values)
  return sql, values
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
// an array of values that correspond to those columns. For each field in the map:
//   * The field name is converted to lower case.
//   * If the field is a nil pointer, the value NULL is used.
func columnsUpdateStringAndVals(mods map[string]interface{}) (string, []interface{}) {
  var keys []string
  var vals []interface{}
  for k, v := range mods {
    if strings.HasPrefix(k, ".") {
      k = strings.TrimPrefix(k, ".")
    }
    // By default, we use all lowercase names for database columns.
    k = strings.ToLower(k)
    if v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) {
      keys = append(keys, k + " = NULL")
    } else {
      keys = append(keys, k + " = ?")
      vals = append(vals, v)
    }
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

// StdDeleteSql generates an SQL DELETE statement.
func stdDeleteByIDSql(tableName string) string {
  sql := "delete from " + tableName + " where id=?;"
  return sql
}
