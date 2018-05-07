package strsql

import (
  "database/sql"
  "fmt"
  "reflect"
  "testing"

  _ "github.com/mattn/go-sqlite3"
)

type eTestRow struct {
  n int;
  s string;
}

func collectETestRows(db *sql.DB, query string) ([]*eTestRow, error) {
  rows := make([]*eTestRow, 0)
  row := &eTestRow{}
  targets := []interface{}{
    &row.n,
    &row.s,
  }
  collector := func() {
    rowCopy := eTestRow(*row)
    rows = append(rows, &rowCopy)
  }
  err := QueryAndCollect(db, query, targets, collector)
  return rows, err
}

func setupAndCollectETestRows(setup, query string) ([]*eTestRow, error) {
  db, err := sql.Open("sqlite3", ":memory:")
  if err != nil {
    return nil, fmt.Errorf("error opening test database: %v", err)
  }
  defer db.Close()

  if err := ExecMulti(db,setup); err != nil {
    return nil, fmt.Errorf("error calling ExecMulti: %v", err)
  }

  return collectETestRows(db, query)
}

func TestExecMulti(t *testing.T) {
  setup := `
CREATE table test(n int, s string);
INSERT into test(n, s) values(1, 'a'), (2, 'b'), (3, 'c');
`
  query := "SELECT n, s from test order by n;"
  expectedResult := []*eTestRow{
    &eTestRow{1, "a"},
    &eTestRow{2, "b"},
    &eTestRow{3, "c"},
  }

  rows, err := setupAndCollectETestRows(setup, query)
  if err != nil {
    t.Fatalf("Error collecting rows: %v", err)
  }

  if got, want := len(rows), 3; got != want {
    t.Fatalf("Wrong number of rows, got %d, want %d", got, want)
  }
  if got, want := rows, expectedResult; !reflect.DeepEqual(got, want) {
    t.Errorf("Results array, got %v, want %v", got, want)
  }
}

func TestComments(t *testing.T) {
  setup := `
CREATE table test(n int, s string);
# This is a comment
INSERT into test(n, s)
# another comment
  values(1, 'a');
`
  query := "SELECT n, s from test order by n;"
  expectedResult := []*eTestRow{
    &eTestRow{1, "a"},
  }

  rows, err := setupAndCollectETestRows(setup, query)
  if err != nil {
    t.Fatalf("Error collecting rows: %v", err)
  }

  if got, want := len(rows), 1; got != want {
    t.Fatalf("Wrong number of rows, got %d, want %d", got, want)
  }
  if got, want := rows, expectedResult; !reflect.DeepEqual(got, want) {
    t.Errorf("Results array, got %v, want %v", got, want)
  }
}
