package structsql

import (
  "reflect"
  "strings"
)

type ColumnInfo struct {
  Name string
  Type string
  Required bool
  IsForeignKey bool
  FKTable string
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
    }
    columnInfos[i].Name = columnName
    columnInfos[i].Type = columnType
  }
  return columnInfos
}
