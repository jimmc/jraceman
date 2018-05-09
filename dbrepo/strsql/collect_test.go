package strsql_test

import (
  "database/sql"
  "reflect"
  "testing"

  "github.com/jimmc/jracemango/dbrepo/dbtesting"
  "github.com/jimmc/jracemango/dbrepo/strsql"
)

type testRow struct {
  n int;
  s string;
}

func setupDatabase(db *sql.DB) error {
  return dbtesting.CreateAndPopulateTestTable(db)
}

func TestCollectNone(t *testing.T) {
  db, err := dbtesting.EmptyDb()
  if err != nil {
    t.Fatalf("Error opening test database: %v", err)
  }
  defer db.Close()

  if err := setupDatabase(db); err != nil {
    t.Fatalf("Error setting up database: %v", err)
  }

  rows := make([]*testRow, 0)
  row := &testRow{}
  targets := []interface{}{
    &row.n,
    &row.s,
  }
  collector := func() {
    rowCopy := testRow(*row)
    rows = append(rows, &rowCopy)
  }
  expectedResult := []*testRow{
    &testRow{1, "a"},
    &testRow{2, "b"},
    &testRow{3, "c"},
  }
  query := "SELECT n, s from test order by n;"
  if err := strsql.QueryAndCollect(db, query, targets, collector); err != nil {
    t.Fatalf("Error collecting rows: %v", err)
  }

  if got, want := len(rows), 3; got != want {
    t.Fatalf("Wrong number of rows, got %d, want %d", got, want)
  }
  if got, want := rows, expectedResult; !reflect.DeepEqual(got, want) {
    t.Errorf("Results array, got %v, want %v", got, want)
  }
}

func TestQueryAndCollectErrors(t *testing.T) {
  db, err := dbtesting.EmptyDb()
  if err != nil {
    t.Fatalf("Error opening test database: %v", err)
  }
  defer db.Close()

  if err := setupDatabase(db); err != nil {
    t.Fatalf("Error setting up database: %v", err)
  }

  // Test the first error return.
  if err := strsql.QueryAndCollect(db, "invalid sql", nil, nil); err == nil {
    t.Errorf("Expected error for invalid sql")
  }

  // Test the second error return (not enough targets).
  if err := strsql.QueryAndCollect(db, "SELECT s from test;", nil, nil); err == nil {
    t.Errorf("Expected error for sql for empty table")
  }
}
