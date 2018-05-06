package strsql

import (
  "database/sql"
  "reflect"
  "testing"

  _ "github.com/mattn/go-sqlite3"
)

type eTestRow struct {
  n int;
  s string;
}

func TestExecMulti(t *testing.T) {
  db, err := sql.Open("sqlite3", ":memory:")
  if err != nil {
    t.Fatalf("Error opening test database: %v", err)
  }
  defer db.Close()

  setupStr := `
CREATE table test(n int, s string);
INSERT into test(n, s) values(1, 'a'), (2, 'b'), (3, 'c');
`
  if err := ExecMulti(db,setupStr); err != nil {
    t.Fatalf("Error calling ExecMulti: %v", err)
  }

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
  expectedResult := []*eTestRow{
    &eTestRow{1, "a"},
    &eTestRow{2, "b"},
    &eTestRow{3, "c"},
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
