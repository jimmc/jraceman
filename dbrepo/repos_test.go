package dbrepo

import (
  "testing"
)

func TestOpenNormal(t *testing.T) {
  dbr, err := Open("ramsql:Test")
  if err != nil {
    t.Errorf("Failed to open ramsql database: %v", err)
  }
  dbr.Close()
}

func TestOpenNoType(t *testing.T) {
  dbr, err := Open(":foo")
  if err == nil {
    dbr.Close()
    t.Errorf("Expected error opening database")
  }

  dbr, err = Open("foo")
  if err == nil {
    dbr.Close()
    t.Errorf("Expected error opening database")
  }
}

func TestCreateTables(t *testing.T) {
  dbr, err := Open("ramsql:Test2")
  if err != nil {
    t.Errorf("Failed to open ramsql database: %v", err)
  }
  defer dbr.Close()
  if err := dbr.CreateTables(); err != nil {
    t.Errorf("Error creating tables: %v", err)
  }
}
