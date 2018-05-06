package structsql

import (
  "database/sql"
  "reflect"
  "testing"

  _ "github.com/mattn/go-sqlite3"
)

type testRow struct {
  n int;
  s string;
}

func setupDatabase(db *sql.DB) error {
  createSql := "CREATE table test(n int, s string);"
  if _, err := db.Exec(createSql); err != nil {
    return err
  }

  insertSql := "INSERT into test(n, s) values(1, 'a'), (2, 'b'), (3, 'c');"
  if _, err := db.Exec(insertSql); err != nil {
    return err
  }

  return nil
}

func TestCollectNone(t *testing.T) {
  db, err := sql.Open("sqlite3", ":memory:")
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
  if err := QueryAndCollect(db, query, targets, collector); err != nil {
    t.Fatalf("Error collecting rows: %v", err)
  }

  if got, want := len(rows), 3; got != want {
    t.Fatalf("Wrong number of rows, got %d, want %d", got, want)
  }
  if got, want := rows, expectedResult; !reflect.DeepEqual(got, want) {
    t.Errorf("Results array, got %v, want %v", got, want)
  }
}
