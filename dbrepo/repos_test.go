package dbrepo

import (
  "bytes"
  "io/ioutil"
  "os"
  "strings"
  "testing"

  _ "github.com/mattn/go-sqlite3"
)

func EmptyRepos() (*Repos, error) {
  return Open("sqlite3::memory:")
}

func TestOpenNormal(t *testing.T) {
  dbr, err := EmptyRepos()
  if err != nil {
    t.Errorf("Failed to open sqlite3 in-memory database: %v", err)
  }
  dbr.Close()

  // Call Close a second time to make sure it does nothing bad.
  dbr.Close()
}

func TestOpenNoType(t *testing.T) {
  dbr, err := Open("foo")
  if err == nil {
    dbr.Close()
    t.Errorf("Expected error opening database")
  }

  dbr, err = Open("badtype:foo")
  if err == nil {
    dbr.Close()
    t.Errorf("Expected error opening database")
  }

}

func TestCreateTables(t *testing.T) {
  dbr, err := EmptyRepos()
  if err != nil {
    t.Errorf("Failed to open sqlite3 in-memory database: %v", err)
  }
  defer dbr.Close()
  if err := dbr.CreateTables(); err != nil {
    t.Errorf("Error creating tables: %v", err)
  }
}

func TestCreateTablesSiteError(t *testing.T) {
  dbr, err := EmptyRepos()
  if err != nil {
    t.Errorf("Failed to open sqlite3 in-memory database: %v", err)
  }
  defer dbr.Close()

  // Create the site table so that we get an error when we try that again.
  if err := dbr.dbSite.CreateTable(); err != nil {
    t.Fatalf("Error creating site table: %v", err)
  }
  err = dbr.CreateTables()
  if err == nil {
    t.Fatalf("Expected error from CreateTables");
  }
  if !strings.Contains(err.Error(), "creating Site table") {
    t.Errorf("Expected error about createing Site table")
  }
}

func TestCreateTablesAreaError(t *testing.T) {
  dbr, err := EmptyRepos()
  if err != nil {
    t.Errorf("Failed to open sqlite3 in-memory database: %v", err)
  }
  defer dbr.Close()

  // Create the site table so that we get an error when we try that again.
  if err := dbr.dbArea.CreateTable(); err != nil {
    t.Fatalf("Error creating area table: %v", err)
  }
  err = dbr.CreateTables()
  if err == nil {
    t.Fatalf("Expected error from CreateTables");
  }
  if !strings.Contains(err.Error(), "creating Area table") {
    t.Errorf("Expected error about createing Area table")
  }
}

func TestImport(t *testing.T) {
  infile, err := os.Open("testdata/import.txt")
  if err != nil {
    t.Fatalf("Error opening import file: %v", err)
  }
  defer infile.Close()
  dbr, err := EmptyRepos()
  if err != nil {
    t.Errorf("Failed to open sqlite3 in-memory database: %v", err)
  }
  defer dbr.Close()
  if err := dbr.CreateTables(); err != nil {
    t.Errorf("Error creating tables: %v", err)
  }

  inserted, updated, unchanged, err := dbr.Import(infile)
  if err != nil {
    t.Errorf("Error importing file: %v", err)
  }
  if got, want := inserted, 5; got != want {
    t.Errorf("Inserted count, got %d, want %d", got, want)
  }
  if got, want := updated, 0; got != want {
    t.Errorf("Updated count, got %d, want %d", got, want)
  }
  if got, want := unchanged, 0; got != want {
    t.Errorf("Unchanged count, got %d, want %d", got, want)
  }

  infile2, err := os.Open("testdata/import2.txt")
  if err != nil {
    t.Fatalf("Error opening import file (2): %v", err)
  }
  defer infile2.Close()
  inserted, updated, unchanged, err = dbr.Import(infile2)
  if err != nil {
    t.Errorf("Error importing file: %v", err)
  }
  if got, want := inserted, 1; got != want {
    t.Errorf("Inserted count, got %d, want %d", got, want)
  }
  if got, want := updated, 2; got != want {
    t.Errorf("Updated count, got %d, want %d", got, want)
  }
  if got, want := unchanged, 3; got != want {
    t.Errorf("Unchanged count, got %d, want %d", got, want)
  }
}

func TestExport(t *testing.T) {
  infile, err := os.Open("testdata/import.txt")
  if err != nil {
    t.Fatalf("Error opening import file: %v", err)
  }
  defer infile.Close()
  dbr, err := EmptyRepos()
  if err != nil {
    t.Errorf("Failed to open sqlite3 in-memory database: %v", err)
  }
  defer dbr.Close()
  if err := dbr.CreateTables(); err != nil {
    t.Errorf("Error creating tables: %v", err)
  }

  _, _, _, err = dbr.Import(infile)
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

  outcontent, err := ioutil.ReadFile(outfilename)
  if err != nil {
    t.Fatalf("Error reading back output file %s: %v", outfilename, err)
  }
  goldencontent, err := ioutil.ReadFile(goldenfilename)
  if err != nil {
    t.Fatalf("Error reading golden file %s: %v", goldenfilename, err)
  }
  if !bytes.Equal(outcontent, goldencontent) {
    t.Errorf("Outfile %s does not match golden file %s", outfilename, goldenfilename)
  }
}
