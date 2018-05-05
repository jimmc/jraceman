package ixport

import (
  "reflect"
  "strings"
  "testing"
)

// TestRowRepo implements the RowRepo interface with methods that capture
// the data so that we can make sure it is correct.
type testRowRepo struct {
  table string
  columns []string
  values []interface{}
  id string
  readCount int
  insertCount int
  updateCount int
}

func (tr *testRowRepo) Read(table string, columns []string, ID string) ([]interface{}, error) {
  tr.table = table
  tr.columns = columns
  tr.id = ID
  tr.readCount++
  var values []interface{}
  switch ID {
  case "A1":
    values = nil        // This is how we indicate we don't have this row
  case "A2":
    values = []interface{}{"A2","abc",true,nil,456}
  case "A3":
    values = []interface{}{"A3","def",false,nil,456}
  default:
    values = nil
  }
  return values, nil
}

func (tr *testRowRepo) Insert(table string, columns[]string, values []interface{}, ID string) error {
  tr.table = table
  tr.columns = columns
  tr.values = values
  tr.id = ID
  tr.insertCount++
  return nil
}

func (tr *testRowRepo) Update(table string, columns[]string, values []interface{}, ID string) error {
  tr.table = table
  tr.columns = columns
  tr.values = values
  tr.id = ID
  tr.updateCount++
  return nil
}

func TestImportDataLineNewRow(t *testing.T) {
  trr := &testRowRepo{}
  im := NewImporter(trr)
  im.tableName = "xtable"
  im.columnNames = []string{"id","scol","bcol","ncol","icol"}
  err := im.importDataLine(`"A1","abc",true,null,456`)
  if err != nil {
    t.Fatalf("Importing data line: %v", err)
  }
  if got, want := trr.table, "xtable"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }
  if got, want := trr.readCount, 1; got != want {
    t.Errorf("ReadCount, got %v, want %v", got, want)
  }
  if got, want := trr.insertCount, 1; got != want {
    t.Errorf("InsertCount, got %v, want %v", got, want)
  }
  if got, want := trr.updateCount, 0; got != want {
    t.Errorf("UpdateCount, got %v, want %v", got, want)
  }
  if got, want := trr.columns, im.columnNames; !reflect.DeepEqual(got, want) {
    t.Errorf("Column names, got %v, want %v", got, want)
  }
}

func TestImportDataLineNoChangeRow(t *testing.T) {
  trr := &testRowRepo{}
  im := NewImporter(trr)
  im.tableName = "xtable"
  im.columnNames = []string{"id","scol","bcol","ncol","icol"}
  err := im.importDataLine(`"A2","abc",true,null,456`)
  if err != nil {
    t.Fatalf("Importing data line: %v", err)
  }
  if got, want := trr.table, "xtable"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }
  if got, want := trr.readCount, 1; got != want {
    t.Errorf("ReadCount, got %v, want %v", got, want)
  }
  if got, want := trr.insertCount, 0; got != want {
    t.Errorf("InsertCount, got %v, want %v", got, want)
  }
  if got, want := trr.updateCount, 0; got != want {
    t.Errorf("UpdateCount, got %v, want %v", got, want)
  }
  if got, want := trr.columns, im.columnNames; !reflect.DeepEqual(got, want) {
    t.Errorf("Column names, got %v, want %v", got, want)
  }
}

func TestImportDataLineUpdateRow(t *testing.T) {
  trr := &testRowRepo{}
  im := NewImporter(trr)
  im.tableName = "xtable"
  im.columnNames = []string{"id","scol","bcol","ncol","icol"}
  if err := im.importDataLine(`"A3","abc",true,null,456`); err != nil {
    t.Fatalf("Importing data line: %v", err)
  }
  if got, want := trr.table, "xtable"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }
  if got, want := trr.readCount, 1; got != want {
    t.Errorf("ReadCount, got %v, want %v", got, want)
  }
  if got, want := trr.insertCount, 0; got != want {
    t.Errorf("InsertCount, got %v, want %v", got, want)
  }
  if got, want := trr.updateCount, 1; got != want {
    t.Errorf("UpdateCount, got %v, want %v", got, want)
  }
  if got, want := trr.columns, []string{"scol", "bcol"}; !reflect.DeepEqual(got, want) {
    t.Errorf("Column names, got %v, want %v", got, want)
  }
}

func TestImportModeLine(t *testing.T) {
  trr := &testRowRepo{}
  im := NewImporter(trr)

  if err := im.importModeLine("!table foo"); err != nil {
    t.Fatalf("Importing mode line (table): %v", err)
  }
  if got, want := im.tableName, "foo"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }

  if err := im.importModeLine(`!columns "id","col2","col3"`); err != nil {
    t.Fatalf("Importing mode line (colmns): %v", err)
  }
  if got, want := im.columnNames, []string{"id", "col2", "col3"}; !reflect.DeepEqual(got, want) {
    t.Errorf("Column names, got %v, want %v", got, want)
  }

  if err := im.importModeLine(`!columns "col1","col2"`); err == nil {
    t.Errorf("Importing mode line (colmns): id col should be required")
  }
}

func TestImportIgnoreLine(t *testing.T) {
  trr := &testRowRepo{}
  im := NewImporter(trr)

  if err := im.importLine(""); err != nil {
    t.Errorf("Importing empty line: %v", err)
  }
  if err := im.importLine("   "); err != nil {
    t.Errorf("Importing spaces line: %v", err)
  }
  if err := im.importLine("# This is a comment line"); err != nil {
    t.Errorf("Importing comment line: %v", err)
  }
  if err := im.importLine("  # This is a comment line"); err != nil {
    t.Errorf("Importing spaces+comment line: %v", err)
  }
}

func TestImportLine(t *testing.T) {
  trr := &testRowRepo{}
  im := NewImporter(trr)

  if err := im.importLine("  ! table foo "); err != nil {
    t.Errorf("Importing table line: %v", err)
  }
  if got, want := im.tableName, "foo"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }
  if err := im.importLine(`  ! columns "id","scol","bcol","ncol","icol" `); err != nil {
    t.Errorf("Importing columns line: %v", err)
  }
  if err := im.importDataLine(`"A1","abc",true,null,456`); err != nil {
    t.Fatalf("Data line: %v", err)
  }
  if got, want := trr.insertCount, 1; got != want {
    t.Errorf("InsertCount, got %v, want %v", got, want)
  }
}

func TestImport(t *testing.T) {
  trr := &testRowRepo{}
  im := NewImporter(trr)
  source := `#!jraceman -import
!exportVersion 1
!appInfo JRaceman v1.1.6
!type database

!table thing
!columns "id","name","meters"
"T1","widget",2
"T2","gadget",4
`
  if err := im.Import(strings.NewReader(source)); err != nil {
    t.Fatalf("Importing string: %v", err)
  }
  if got, want := im.tableName, "thing"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }
  if got, want := trr.readCount, 2; got != want {
    t.Errorf("ReadCount, got %v, want %v", got, want)
  }
  if got, want := trr.insertCount, 2; got != want {
    t.Errorf("InsertCount, got %v, want %v", got, want)
  }
}
