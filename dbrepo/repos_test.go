package dbrepo_test

import (
  "os"
  "strings"
  "testing"

  "github.com/jimmc/jraceman/dbrepo"
  dbtest "github.com/jimmc/jraceman/dbrepo/test"

  goldenbase "github.com/jimmc/golden/base"
)

func TestOpenNormal(t *testing.T) {
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    t.Errorf("Failed to open test database: %v", err)
  }
  dbr.Close()

  // Call Close a second time to make sure it does nothing bad.
  dbr.Close()
}

func TestOpenNoType(t *testing.T) {
  dbr, err := dbrepo.Open("foo")
  if err == nil {
    dbr.Close()
    t.Errorf("Expected error opening database")
  }

  dbr, err = dbrepo.Open("badtype:foo")
  if err == nil {
    dbr.Close()
    t.Errorf("Expected error opening database")
  }

}

func TestCreateTables(t *testing.T) {
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    t.Errorf("Failed to open test database: %v", err)
  }
  defer dbr.Close()
  if err := dbr.CreateTables(); err != nil {
    t.Errorf("Error creating tables: %v", err)
  }
}

func TestCreateTablesSiteError(t *testing.T) {
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    t.Errorf("Failed to open test database: %v", err)
  }
  defer dbr.Close()

  // Create the site table so that we get an error when we try that again.
  if err := dbr.Site().(*dbrepo.DBSiteRepo).CreateTable(); err != nil {
    t.Fatalf("Error creating site table: %v", err)
  }
  err = dbr.CreateTables()
  if err == nil {
    t.Fatalf("Expected error from CreateTables");
  }
  if !strings.Contains(err.Error(), "table site already exists") {
    t.Errorf("Expected error about creating Site table, got: %v", err)
  }
}

func TestCreateTablesAreaError(t *testing.T) {
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    t.Errorf("Failed to open test database: %v", err)
  }
  defer dbr.Close()

  // Create the site table so that we get an error when we try that again.
  if err := dbr.Area().(*dbrepo.DBAreaRepo).CreateTable(); err != nil {
    t.Fatalf("Error creating area table: %v", err)
  }
  err = dbr.CreateTables()
  if err == nil {
    t.Fatalf("Expected error from CreateTables");
  }
  if !strings.Contains(err.Error(), "table area already exists") {
    t.Errorf("Expected error about creating Area table, got: %v", err)
  }
}

func TestImport(t *testing.T) {
  filename := "testdata/import.txt"
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    t.Errorf("Failed to open test database: %v", err)
  }
  defer dbr.Close()
  if err := dbr.CreateTables(); err != nil {
    t.Errorf("Error creating tables: %v", err)
  }

  counts, err := dbr.ImportFile(filename)
  if err != nil {
    t.Errorf("Error importing file: %v", err)
  }
  if got, want := counts.Inserted(), 5; got != want {
    t.Errorf("Inserted count, got %d, want %d", got, want)
  }
  if got, want := counts.Updated(), 0; got != want {
    t.Errorf("Updated count, got %d, want %d", got, want)
  }
  if got, want := counts.Unchanged(), 0; got != want {
    t.Errorf("Unchanged count, got %d, want %d", got, want)
  }

  filename2 := "testdata/import2.txt"
  counts, err = dbr.ImportFile(filename2)
  if err != nil {
    t.Errorf("Error importing file: %v", err)
  }
  if got, want := counts.Inserted(), 1; got != want {
    t.Errorf("Inserted count, got %d, want %d", got, want)
  }
  if got, want := counts.Updated(), 2; got != want {
    t.Errorf("Updated count, got %d, want %d", got, want)
  }
  if got, want := counts.Unchanged(), 3; got != want {
    t.Errorf("Unchanged count, got %d, want %d", got, want)
  }
}

func TestExport(t *testing.T) {
  filename := "testdata/import.txt"
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    t.Errorf("Failed to open test database: %v", err)
  }
  defer dbr.Close()
  if err := dbr.CreateTables(); err != nil {
    t.Errorf("Error creating tables: %v", err)
  }

  _, err = dbr.ImportFile(filename)
  if err != nil {
    t.Errorf("Error importing file: %v", err)
  }

  outfilename := "testdata/export.out"
  goldenfilename := "testdata/export.golden"
  os.Remove(outfilename)
  outfile, err := os.Create(outfilename)
  if err != nil {
    t.Fatalf("Error creating export file: %v", err)
  }
  defer outfile.Close()
  // We leave the outfile there so we can rename it to
  // the golden file if we want.

  err = dbr.Export(outfile)
  if err != nil {
    t.Fatalf("Error exporting: %v", err)
  }

  if err := goldenbase.CompareOutToGolden(outfilename, goldenfilename); err != nil {
    t.Error(err.Error())
  }
}
