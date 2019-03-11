package structsql

import (
  "log"
  "reflect"
  "strings"
)

// ColumnInfo describes the attributes of a database column.
type ColumnInfo struct {
  Name string
  Type string
  Required bool
  IsForeignKey bool
  FKTable string
  HasDefault bool
  DefaultAsString string
}

// ColumnInfosDiff describes how to change from one []ColumnInfo
// to another []ColumnInfo by matching up the Name fields.
type ColumnInfosDiff struct {
  Add []ColumnInfo
  Change [][2]ColumnInfo
  Remove []ColumnInfo
}

// ColumnNames generates a list of column names
// based on the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case.
func ColumnNames(entity interface{}) []string {
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

// ColumnInfo generates a list of column names and other info
// from the given struct.
func ColumnInfos(entity interface{}) []ColumnInfo {
  val := reflect.Indirect(reflect.ValueOf(entity))
  typ := val.Type()
  numFields := typ.NumField()
  columnInfos := make([]ColumnInfo, numFields)
  for i := 0; i < numFields; i++ {
    field := typ.Field(i)
    columnName := strings.ToLower(field.Name)
    goTypeName := field.Type.String()     // string, *string, int
    isPointer := strings.HasPrefix(goTypeName, "*")
    isForeignKey := strings.HasSuffix(field.Name, "ID")
    goTypeName = strings.TrimPrefix(goTypeName, "*")
    columnType := goTypeName            // TODO - convert as required
    if columnName == "id" {
      columnInfos[i].Required = true
    } else {
      columnInfos[i].Required = !isPointer
      if isForeignKey {
        columnInfos[i].IsForeignKey = isForeignKey
        columnInfos[i].FKTable = strings.TrimSuffix(columnName, "id")
      }
      if columnInfos[i].Required {
        // For non-id required columns, assume the default is the zero value.
        switch columnType {
        // Only put in default for bool.
        case "bool":
          columnInfos[i].HasDefault = true
          columnInfos[i].DefaultAsString = "false"
        // Allow higher-level code to put in defaults for specified fields.
        // case "float32":
        //  columnInfos[i].HasDefault = true
        //  columnInfos[i].DefaultAsString = "0.0"
        // case "int":
        //  columnInfos[i].HasDefault = true
        //  columnInfos[i].DefaultAsString = "0"
        // case "string":
        //  columnInfos[i].HasDefault = true
        //  columnInfos[i].DefaultAsString = "''"
        default: log.Printf("Unknown columnType %q for %s", columnType, columnName)
        }
      }
    }
    columnInfos[i].Name = columnName
    columnInfos[i].Type = columnType
  }
  return columnInfos
}

// DiffColumnInfos compares two arrays of ColumnInfo and returns a struct
// with info on what needs to be done to go from got to want.
func DiffColumnInfos(got, want []ColumnInfo) ColumnInfosDiff {
  gotMap := ColumnInfosToMap(got)
  wantMap := ColumnInfosToMap(want)
  add := make([]ColumnInfo, 0)
  change := make([][2]ColumnInfo, 0)
  remove := make([]ColumnInfo, 0)
  for _, g := range got {
    w := wantMap[g.Name]
    if w == nil {
      // We got it, but we don't want it
      remove = append(remove,  g)
    } else {
      // We got it, and we want it...
      if ColumnInfosDiffer(&g, w) {
        // ... but what we have is different from what we want.
        change = append(change, [2]ColumnInfo{g, *w})
      } else {
        // ... and it's the same, so no need to do anything.
      }
    }
  }
  for _, w := range want {
    g := gotMap[w.Name]
    if g == nil {
      // We want it, but we don't got it
      add = append(add, w)
    }
  }
  return ColumnInfosDiff{
    Add: add,
    Change: change,
    Remove: remove,
  }
}

// ColumnInfosDiffer returns true if any of the fields of
// the two args differ.
func ColumnInfosDiffer(g, w *ColumnInfo) bool {
  if g.Name != w.Name {
    return true
  }
  if g.Type != w.Type {
    return true
  }
  if g.Required != w.Required {
    return true
  }
  if g.IsForeignKey != w.IsForeignKey {
    return true
  }
  if g.FKTable != w.FKTable {
    return true
  }
  return false
}

// ColumnInfosToMap takes an array of ColumnInfo and returns a map
// of those structures using the Name field as the key.
func ColumnInfosToMap(cols []ColumnInfo) map[string]*ColumnInfo {
  m := make(map[string]*ColumnInfo, len(cols))
  for i, c := range cols {
    m[c.Name] = &cols[i]
  }
  return m
}
