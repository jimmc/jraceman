package ixport_test

import (
  "database/sql"
  "io"
  "os"
  "testing"

  "github.com/jimmc/jracemango/dbrepo/dbtesting"
  "github.com/jimmc/jracemango/dbrepo/ixport"
)

func TestExportHeader(t *testing.T) {
  db, err := dbtesting.EmptyDb()
  if err != nil {
    t.Fatalf("Error creating test database: %v", err)
  }
  defer db.Close()

  outfilename := "testdata/exportheader.out"
  goldenfilename := "testdata/exportheader.golden"
  os.Remove(outfilename)
  outfile, err := os.Create(outfilename)
  if err != nil {
    t.Fatalf("Error creating output file: %v", outfile)
  }

  e := ixport.NewExporter(db)
  e.ExportHeader(outfile)

  outfile.Close()
  if err := dbtesting.CompareOutToGolden(outfilename, goldenfilename); err != nil {
    t.Error(err.Error())
  }
}

type testEntity struct {
  N int
  S string
  S2 *string
}

func TestExportTable(t *testing.T) {
  callback := func(db *sql.DB, outfile io.Writer) error {
    e := ixport.NewExporter(db)
    return e.ExportTableFromStruct(outfile, "test", &testEntity{})
  }
  if err := dbtesting.FromSetupToGolden("exporttable", callback); err != nil {
    t.Error(err.Error())
  }
}

func TestExportEmptyTable(t *testing.T) {
  callback := func(db *sql.DB, outfile io.Writer) error {
    e := ixport.NewExporter(db)
    return e.ExportTableFromStruct(outfile, "test", &testEntity{})
  }
  if err := dbtesting.FromSetupToGolden("exportemptytable", callback); err != nil {
    t.Error(err.Error())
  }
}
