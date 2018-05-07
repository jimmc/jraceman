package dbrepo

import (
  "testing"

  _ "github.com/mattn/go-sqlite3"
)

func TestOpenNormal(t *testing.T) {
  dbr, err := Open("sqlite3::memory:")
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
  dbr, err := Open("sqlite3::memory:")
  if err != nil {
    t.Errorf("Failed to open sqlite3 in-memory database: %v", err)
  }
  defer dbr.Close()
  if err := dbr.CreateTables(); err != nil {
    t.Errorf("Error creating tables: %v", err)
  }
}
