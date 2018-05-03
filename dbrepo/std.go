package dbrepo

import (
  "database/sql"
  "fmt"
  "log"
  "reflect"
  "sort"
  "strconv"
  "strings"
)

// Placeholder is ? for MySQL,$N for PostgreSQL,
// SQLite uses either of those, Oracle is :param1

// StdColumnNamesFromStruct generates a list of column names
// based on the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
func stdColumnNamesFromStruct(entity interface{}) []string {
  val := reflect.Indirect(reflect.ValueOf(entity))
  typ := val.Type()
  numFields := typ.NumField()
  columnNames := make([]string, numFields)
  for i := 0; i < numFields; i++ {
    field := typ.Field(i)
    columnNames[i] = strings.ToLower(field.Name)
  }
  return columnNames
}

func stdCreateTableFromStruct(db *sql.DB, tableName string, entity interface{}) error {
  sql := stdCreateTableSqlFromStruct(tableName, entity)
  _, err := db.Exec(sql)
  return err
}

// StdCreateTableSqlFromStruct generates an SQL CREATE TABLE command using
// the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case for the column name.
//   * int and string fields are declared as that same type column.
//   * The id field is declared as primary key.
//   * Non-pointer fields are declared as not null.
//   * Field names ending in ID are declared as foreign key references to the
//     id field of a table whose name matches the first part of the field name
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
    isForeignKey := strings.HasSuffix(field.Name, "ID")
    goTypeName = strings.TrimPrefix(goTypeName, "*")
    columnType := goTypeName            // TODO - convert as required
    columnSpec := columnName + " " + columnType
    if columnName == "id" {
      columnSpec = columnSpec + " primary key"
    } else {
      if !isPointer {
        columnSpec = columnSpec + " not null"
      }
      if isForeignKey {
        referenceTable := strings.TrimSuffix(columnName, "id")
        columnSpec = columnSpec + " references " + referenceTable + "(id)"
      }
    }
    columnSpecs[i] = columnSpec
  }
  sql := "CREATE TABLE " + tableName + "(" + strings.Join(columnSpecs, ", ") + ");"
  log.Printf("stdCreateTableSql: %v\n", sql)
  return sql
}

// StdSelectSqlFromStruct generates an SQL QUERY statement
// using the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
// There is no semicolon at the end of the sql, so the caller can append more
// sql to it such as a where clause.
func stdSelectSqlFromStruct(tableName string, entity interface{}) (string, []interface{}) {
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

// StdFindByIDSqlFromStruct generates an SQL QUERY statement,
// with a WHERE clause limiting it to a matching id,
// using the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
func stdFindByIDSqlFromStruct(tableName string, entity interface{}) (string, []interface{}) {
  sql, targets := stdSelectSqlFromStruct(tableName, entity)
  sql = sql + " where id=?;"
  log.Printf("stdFindByIDSql: %v\n", sql)
  return sql, targets
}

// StdListSqlFromStruct generates an SQL QUERY statement,
// with OFFSET and LIMIT clauses,
// using the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
func stdListSqlFromStruct(tableName string, entity interface{}, offset, limit int) (string, []interface{}) {
  sql, targets := stdSelectSqlFromStruct(tableName, entity)
  if limit != 0 {
    sql = sql + " limit " + strconv.Itoa(limit)
  }
  if offset != 0 {
    sql = sql + " offset " + strconv.Itoa(offset)
  }
  log.Printf("stdListSql: %v\n", sql)
  return sql, targets
}

func stdInsertFromStruct(db *sql.DB, tableName string, entity interface{}, ID string) error {
  sql, values := stdInsertSqlFromStruct(tableName, entity)
  res, err := db.Exec(sql, values...)
  return requireOneResult(res, err, "Inserted", tableName, ID)
}

// StdInsertSqlFromStruct generates an SQL INSERT statement using
// the fields of the given struct. For each field in the struct:
//   * If the field is a nil pointer, it is ignored.
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

// StdQueryAndCollect issues a Query for the given sql, then interates through
// the returned rows. For each row, it retrieves the data into targets, then
// calls the collect function. The assumption is that the targets store the
// results into data that is accessible to the collect function.
func stdQueryAndCollect(db *sql.DB, sql string, targets []interface{}, collect func()) error {
  rows, err := db.Query(sql)
  if err != nil {
    return err
  }
  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(targets...)
    if err != nil {
      return err
    }
    collect()
  }
  return rows.Err()
}

// We need only the RowsAffected function from the database/sql.Result interface.
type sqlRowsAffected interface {
  RowsAffected() (int64, error)
}

// RequireOneResult gets the result of sql.Stmt.Exec and verifies that it
// affected exactly one row, which should be the case for operations that
// use the entity ID. The res argument is typically a sql.Result.
// The action string is used in error messages, and
// should be capitalized and past tense, such as "Deleted".
// The entityType should the name of the entity, such as "site".
func requireOneResult(res sqlRowsAffected, err error, action, entityType, ID string) error {
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

// StdUpdateByID updates a record by ID.
func stdUpdateByID(db *sql.DB, tableName string, mods map[string]interface{}, ID string) error {
  sql, vals := modsToSql(tableName, mods, ID)
  res, err := db.Exec(sql, vals...)
  return requireOneResult(res, err, "Updated", tableName, ID)
}

// ColumnsUpdateStringAndValues generates a string for the column-and-values portion
// of an SQL update statement, in the form "col1 = ?, col2 = ?", and also returns
// an array of values that correspond to those columns. For each field in the map:
//   * The field name is converted to lower case.
//   * If the field is a nil pointer, the value NULL is used.
func columnsUpdateStringAndVals(mods map[string]interface{}) (string, []interface{}) {
  // We get the keys and sort them so that we have a determinisitic ordering.
  allkeys := make([]string, len(mods))
  i := 0
  for k := range mods {
    allkeys[i] = k
    i++
  }
  sort.Strings(allkeys)
  var keys []string
  var vals []interface{}
  for _, k := range allkeys {
    v := mods[k]
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
  updateSql := "update " + table + " set " + kvString + " where id = ?;"
  vals = append(vals, ID)
  log.Printf("updateSql = %q, vals = %v", updateSql, vals)
  return updateSql, vals
}

// StdDelete deletes a record by ID.
func stdDeleteByID(db *sql.DB, tableName, ID string) error {
  sql := stdDeleteByIDSql(tableName)
  res, err := db.Exec(sql, ID)
  return requireOneResult(res, err, "Deleted", "site", ID)
}

// StdDeleteSql generates an SQL DELETE statement.
func stdDeleteByIDSql(tableName string) string {
  sql := "delete from " + tableName + " where id=?;"
  return sql
}
