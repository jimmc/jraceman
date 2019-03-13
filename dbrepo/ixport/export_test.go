package ixport_test

import (
  "database/sql"
  "io"
  "os"
  "testing"

  dbtest "github.com/jimmc/jracemango/dbrepo/test"
  "github.com/jimmc/jracemango/dbrepo/ixport"
)

func TestExportHeader(t *testing.T) {
  db, err := dbtest.EmptyDb()
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
  if err := dbtest.CompareOutToGolden(outfilename, goldenfilename); err != nil {
    t.Error(err.Error())
  }
}

type testEntity struct {
  N int
  S string
  S2 *string
  B bool
  F float32
}

func TestExportTable(t *testing.T) {
  callback := func(db *sql.DB, outfile io.Writer) error {
    e := ixport.NewExporter(db)
    return e.ExportTableFromStruct(outfile, "test", &testEntity{})
  }
  if err := dbtest.FromSetupToGolden("exporttable", callback); err != nil {
    t.Error(err.Error())
  }
}

func TestExportEmptyTable(t *testing.T) {
  callback := func(db *sql.DB, outfile io.Writer) error {
    e := ixport.NewExporter(db)
    return e.ExportTableFromStruct(outfile, "test", &testEntity{})
  }
  if err := dbtest.FromSetupToGolden("exportemptytable", callback); err != nil {
    t.Error(err.Error())
  }
}
