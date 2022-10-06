package ixport_test

import (
  "reflect"
  "strings"
  "testing"

  "github.com/jimmc/jracemango/dbrepo/ixport"
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

func (tr *testRowRepo) ReadByKey(table string, columns []string, keyName, key string) ([]interface{}, error) {
  tr.table = table
  tr.columns = columns
  tr.id = key
  tr.readCount++
  var values []interface{}
  switch key {
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

func (tr *testRowRepo) UpdateByKey(table string, columns[]string, values []interface{}, keyName, key string) error {
  tr.table = table
  tr.columns = columns
  tr.values = values
  tr.id = key
  tr.updateCount++
  return nil
}

func TestImportDataLineNewRow(t *testing.T) {
  trr := &testRowRepo{}
  im := ixport.NewImporter(trr)
  if err := im.ImportLine("!table xtable"); err != nil {
    t.Fatalf("Setting table: %v", err)
  }
  if err := im.ImportLine(`!columns "id","scol","bcol","ncol","icol"`); err != nil {
    t.Fatalf("Setting columns: %v", err)
  }
  if err := im.ImportLine(`"A1","abc",true,null,456`); err != nil {
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
  if got, want := trr.columns, im.ColumnNames(); !reflect.DeepEqual(got, want) {
    t.Errorf("Column names, got %v, want %v", got, want)
  }
}

func TestImportDataLineNoChangeRow(t *testing.T) {
  trr := &testRowRepo{}
  im := ixport.NewImporter(trr)
  if err := im.ImportLine("!table xtable"); err != nil {
    t.Fatalf("Setting table: %v", err)
  }
  if err := im.ImportLine(`!columns "id","scol","bcol","ncol","icol"`); err != nil {
    t.Fatalf("Setting columns: %v", err)
  }
  if err := im.ImportLine(`"A2","abc",true,null,456`); err != nil {
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
  if got, want := trr.columns, im.ColumnNames(); !reflect.DeepEqual(got, want) {
    t.Errorf("Column names, got %v, want %v", got, want)
  }
}

func TestImportDataLineUpdateRow(t *testing.T) {
  trr := &testRowRepo{}
  im := ixport.NewImporter(trr)
  if err := im.ImportLine("!table xtable"); err != nil {
    t.Fatalf("Setting table: %v", err)
  }
  if err := im.ImportLine(`!columns "id","scol","bcol","ncol","icol"`); err != nil {
    t.Fatalf("Setting columns: %v", err)
  }
  if err := im.ImportLine(`"A3","abc",true,null,456`); err != nil {
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
  im := ixport.NewImporter(trr)

  if err := im.ImportLine("!table foo"); err != nil {
    t.Fatalf("Importing mode line (table): %v", err)
  }
  if got, want := im.TableName(), "foo"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }

  if err := im.ImportLine(`!columns "id","col2","col3"`); err != nil {
    t.Fatalf("Importing mode line (colmns): %v", err)
  }
  if got, want := im.ColumnNames(), []string{"id", "col2", "col3"}; !reflect.DeepEqual(got, want) {
    t.Errorf("Column names, got %v, want %v", got, want)
  }

  if err := im.ImportLine(`!columns "col1","col2"`); err == nil {
    t.Errorf("Importing mode line (colmns): id col should be required")
  }
}

func TestImportIgnoreLine(t *testing.T) {
  trr := &testRowRepo{}
  im := ixport.NewImporter(trr)

  if err := im.ImportLine(""); err != nil {
    t.Errorf("Importing empty line: %v", err)
  }
  if err := im.ImportLine("   "); err != nil {
    t.Errorf("Importing spaces line: %v", err)
  }
  if err := im.ImportLine("# This is a comment line"); err != nil {
    t.Errorf("Importing comment line: %v", err)
  }
  if err := im.ImportLine("  # This is a comment line"); err != nil {
    t.Errorf("Importing spaces+comment line: %v", err)
  }
}

func TestImportLine(t *testing.T) {
  trr := &testRowRepo{}
  im := ixport.NewImporter(trr)

  if err := im.ImportLine("  ! table foo "); err != nil {
    t.Errorf("Importing table line: %v", err)
  }
  if got, want := im.TableName(), "foo"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }
  if err := im.ImportLine(`  ! columns "id","scol","bcol","ncol","icol" `); err != nil {
    t.Errorf("Importing columns line: %v", err)
  }
  if err := im.ImportLine(`"A1","abc",true,null,456`); err != nil {
    t.Fatalf("Data line: %v", err)
  }
  if got, want := trr.insertCount, 1; got != want {
    t.Errorf("InsertCount, got %v, want %v", got, want)
  }
}

func TestImportv2(t *testing.T) {
  trr := &testRowRepo{}
  im := ixport.NewImporter(trr)
  source := `#!jraceman -import
!exportVersion 2
!appInfo JRaceman v2.0.0
!type database

!table thing
!columns "id","name","meters"
"T1","widget",2
"T2","gadget",4
`
  if err := im.ImportReader(strings.NewReader(source)); err != nil {
    t.Fatalf("Importing string: %v", err)
  }
  if got, want := im.TableName(), "thing"; got != want {
    t.Errorf("Table name, got %v, want %v", got, want)
  }
  if got, want := trr.readCount, 2; got != want {
    t.Errorf("ReadCount, got %v, want %v", got, want)
  }
  if got, want := trr.insertCount, 2; got != want {
    t.Errorf("InsertCount, got %v, want %v", got, want)
  }
  inserted := im.Counts().Inserted()
  updated := im.Counts().Updated()
  unchanged := im.Counts().Unchanged()
  if got, want := inserted, 2; got != want {
    t.Errorf("Import InsertedCount, got %d, want %d", got, want)
  }
  if got, want := updated, 0; got != want {
    t.Errorf("Import UpdatedCount, got %d, want %d", got, want)
  }
  if got, want := unchanged, 0; got != want {
    t.Errorf("Import UnchangedCount, got %d, want %d", got, want)
  }
}
