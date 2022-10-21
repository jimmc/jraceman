package structsql_test

import (
  "reflect"
  "testing"

  "github.com/jimmc/jraceman/dbrepo/structsql"
)

func TestColumnNames(t *testing.T) {
  columnNames := structsql.ColumnNames(foo)
  expectedNames := []string{"id", "num", "required", "optional", "barid"}
  if got, want := columnNames, expectedNames; !reflect.DeepEqual(got, want) {
    t.Errorf("Column names: got %v, want %v", got, want)
  }
}

func TestColumnInfos(t *testing.T) {
  columnInfos := structsql.ColumnInfos(foo)
  expectedInfos := []structsql.ColumnInfo{
    structsql.ColumnInfo{Name: "id", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "num", Type: "int", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "required", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "optional", Type: "string", Required: false, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "barid", Type: "string", Required: false, IsForeignKey: true, FKTable: "bar"},
  }
  if got, want := columnInfos, expectedInfos; !reflect.DeepEqual(got, want) {
    t.Errorf("ColumnInfos: got %v, want %v", got, want)
  }
}

func TestDiffColumnInfos(t *testing.T) {
  gCols := []structsql.ColumnInfo {
    structsql.ColumnInfo{Name: "id", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "str", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "num", Type: "int", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "required", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "optional", Type: "string", Required: false, IsForeignKey: false, FKTable: ""},
  }
  wCols := []structsql.ColumnInfo {
    structsql.ColumnInfo{Name: "id", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "str", Type: "string", Required: false, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "num", Type: "int", Required: false, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "optional", Type: "string", Required: false, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "barid", Type: "string", Required: false, IsForeignKey: true, FKTable: "bar"},
  }
  expectedDiffs := structsql.ColumnInfosDiff{
    Add: []structsql.ColumnInfo{
      structsql.ColumnInfo{Name: "barid", Type: "string", Required: false, IsForeignKey: true, FKTable: "bar"},
    },
    Change: [][2]structsql.ColumnInfo{
      [2]structsql.ColumnInfo{
        structsql.ColumnInfo{Name: "str", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
        structsql.ColumnInfo{Name: "str", Type: "string", Required: false, IsForeignKey: false, FKTable: ""},
      },
      [2]structsql.ColumnInfo{
        structsql.ColumnInfo{Name: "num", Type: "int", Required: true, IsForeignKey: false, FKTable: ""},
        structsql.ColumnInfo{Name: "num", Type: "int", Required: false, IsForeignKey: false, FKTable: ""},
      },
    },
    Remove: []structsql.ColumnInfo{
      structsql.ColumnInfo{Name: "required", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
    },
  }
  diffs := structsql.DiffColumnInfos(gCols, wCols)
  if got, want := diffs, expectedDiffs; !reflect.DeepEqual(got, want) {
    t.Errorf("ColumnInfoDiffs: got %v, want %v", got, want)
  }
}

func TestColumnInfosDiffer(t *testing.T) {
  a := &structsql.ColumnInfo{Name: "foo", Type: "string", Required: true, IsForeignKey: false, FKTable: ""}
  b := &structsql.ColumnInfo{Name: "foo", Type: "int", Required: true, IsForeignKey: false, FKTable: ""}
  if structsql.ColumnInfosDiffer(a, a) {
    t.Errorf("ColInfoDiffs for self compare should not differ")
  }
  if !structsql.ColumnInfosDiffer(a, b) {
    t.Errorf("ColInfoDiffs for different info should differ")
  }
}

func TestColumnInfosToMap(t *testing.T) {
  a := []structsql.ColumnInfo {
    structsql.ColumnInfo{Name: "id", Type: "string", Required: true, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "num", Type: "int", Required: false, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "optional", Type: "string", Required: false, IsForeignKey: false, FKTable: ""},
    structsql.ColumnInfo{Name: "barid", Type: "string", Required: false, IsForeignKey: true, FKTable: "bar"},
  }
  m := structsql.ColumnInfosToMap(a)
  if got, want := len(m), 4; got != want {
    t.Errorf("Map length: got %d, want %d", got, want)
  }
  if got, want := m["id"], &a[0]; got != want {
    t.Errorf("Map entry for 'id': got %v, want %v", got, want)
  }
}
