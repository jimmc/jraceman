package strsql_test

import (
  "reflect"
  "testing"

  dbtest "github.com/jimmc/jracemango/dbrepo/test"
  "github.com/jimmc/jracemango/dbrepo/strsql"

  goldendb "github.com/jimmc/golden/db"
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

func TestQueryIntHappy(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  n, err := strsql.QueryInt(db, "select 1")
  if err != nil {
    t.Fatalf("QueryInt for constant: %v", err)
  }
  if got, want := n, 1; got != want {
    t.Errorf("QueryInt for constant, got %d, want %d", got, want)
  }

  nval := 1
  n, err = strsql.QueryInt(db, "select count(*) from test where n!=?", nval)
  if err != nil {
    t.Fatalf("QueryInt for count: %v", err)
  }
  if got, want := n, 2; got != want {
    t.Errorf("QueryInt for count, got %d, want %d", got, want)
  }
}

func TestQueryIntError(t *testing.T) {
  db, err := goldendb.EmptyDb()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  _, err = strsql.QueryInt(db, "select x from nosuchtable")
  if err == nil {
    t.Fatalf("Expected error from QueryInt")
  }
}

func TestQueryStringHappy(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  s, err := strsql.QueryString(db, `select "abc"`)
  if err != nil {
    t.Fatalf("QueryString for constant: %v", err)
  }
  if got, want := s, "abc"; got != want {
    t.Errorf("QueryString for constant, got %v, want %v", got, want)
  }

  idval := "T1"
  s, err = strsql.QueryString(db, "select s from test where id=?", idval)
  if err != nil {
    t.Fatalf("QueryString for count: %v", err)
  }
  if got, want := s, "a"; got != want {
    t.Errorf("QueryString for count, got %v, want %v", got, want)
  }
}

func TestQueryStringError(t *testing.T) {
  db, err := goldendb.EmptyDb()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  _, err = strsql.QueryString(db, "select x from nosuchtable")
  if err == nil {
    t.Fatalf("Expected error from QueryString")
  }
}

func TestQueryStringsHappy(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  idval := "T1"
  ss, err := strsql.QueryStrings(db, "select s from test where id!=?", idval)
  if err != nil {
    t.Fatalf("QueryString for count: %v", err)
  }
  if got, want := ss, []string{"b", "c"}; !reflect.DeepEqual(got, want) {
    t.Errorf("QueryString for count, got %v, want %v", got, want)
  }
}
