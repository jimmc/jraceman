package structsql

import (
  "database/sql"
  "reflect"
  "sort"
  "strings"

  "github.com/golang/glog"
)

// UpdateByID updates a record by ID.
func UpdateByID(db *sql.DB, tableName string, mods map[string]interface{}, ID string) error {
  sql, vals := ModsToSql(tableName, mods, ID)
  res, err := db.Exec(sql, vals...)
  return RequireOneResult(res, err, "Updated", tableName, ID)
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
func ModsToSql(table string, mods map[string]interface{}, ID string) (string, []interface{}) {
  glog.V(1).Infof("mods = %v", mods)
  kvString, vals := columnsUpdateStringAndVals(mods)
  updateSql := "update " + table + " set " + kvString + " where id = ?;"
  vals = append(vals, ID)
  glog.V(1).Infof("updateSql = %q, vals = %v", updateSql, vals)
  return updateSql, vals
}
