package structsql_test

import (
  "testing"

  "github.com/jimmc/jraceman/dbrepo/structsql"
  dbtest "github.com/jimmc/jraceman/dbrepo/test"
)

func TestUniqueId(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  if got, want := structsql.UniqueID(db, "test", "T1"), "T4"; got != want {
    t.Errorf("UniqueID: got %v, want %v", got, want)
  }
  if got, want := structsql.UniqueID(db, "test", "X1"), "X1"; got != want {
    t.Errorf("UniqueID: got %v, want %v", got, want)
  }
}

func TestNumberNotIn(t *testing.T) {
  cases := []struct{
    expected int
    nn []int
  } {
    { 1, []int{} },
    { 3, []int{4, 5, 6} },
    { 5, []int{1, 2, 3, 4} },
    { 5, []int{0, 1, 2, 3, 4} },
    { 5, []int{-2, -1, 0, 1, 2, 3, 4} },
    { 5, []int{1, 2, 6, 9} },
    { 4, []int{-4, -2, 1, 2, 5, 7} },
    { 2, []int{-4, -2, 3, 5, 7} },
  }
  for _, c := range cases {
    if got, want := structsql.NumberNotIn(c.nn), c.expected; got != want {
      t.Errorf("NumberNotIn(%v): got %v, want %v", c.nn, got, want)
    }
  }
}
