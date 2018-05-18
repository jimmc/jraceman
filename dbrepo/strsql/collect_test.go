package strsql_test

import (
  "reflect"
  "testing"

  dbtest "github.com/jimmc/jracemango/dbrepo/test"
  "github.com/jimmc/jracemango/dbrepo/strsql"
)

type testRow struct {
  n int;
  s string;
}

func TestCollectNone(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

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
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  // Test the first error return.
  if err := strsql.QueryAndCollect(db, "invalid sql", nil, nil); err == nil {
    t.Errorf("Expected error for invalid sql")
  }

  // Test the second error return (not enough targets).
  if err := strsql.QueryAndCollect(db, "SELECT s from test;", nil, nil); err == nil {
    t.Errorf("Expected error for sql for empty table")
  }
}

func TestQueryStarAndCollect(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  var expectedRowsResult [][]interface{} = make([][]interface{}, 3)
  expectedRowsResult[0] = make([]interface{}, 2)
  expectedRowsResult[1] = make([]interface{}, 2)
  expectedRowsResult[2] = make([]interface{}, 2)
  expectedRowsResult[0][0] = int64(1)
  expectedRowsResult[0][1] = "a"
  expectedRowsResult[1][0] = int64(2)
  expectedRowsResult[1][1] = "b"
  expectedRowsResult[2][0] = int64(3)
  expectedRowsResult[2][1] = "c"

  expectedColumnsResult := []*strsql.ColumnInfo{
    &strsql.ColumnInfo{
      Name: "n",
      Type: "int",
    },
    &strsql.ColumnInfo{
      Name: "s",
      Type: "string",
    },
  }

  query := "SELECT n, s from test order by n;"
  results, err := strsql.QueryStarAndCollect(db, query)
  if err != nil {
    t.Fatalf("Error collecting rows: %v", err)
  }

  if got, want := len(results.Columns), 2; got != want {
    t.Fatalf("Wrong number of columns, got %d, want %d", got, want)
  }
  if got, want := len(results.Rows), 3; got != want {
    t.Fatalf("Wrong number of rows, got %d, want %d", got, want)
  }
  if got, want := results.Columns, expectedColumnsResult; !reflect.DeepEqual(got, want) {
    t.Errorf("Results columns, got %v, want %v", got, want)
  }
  if got, want := results.Rows, expectedRowsResult; !reflect.DeepEqual(got, want) {
    t.Errorf("Results rows, got %v, want %v", got, want)
  }
}
