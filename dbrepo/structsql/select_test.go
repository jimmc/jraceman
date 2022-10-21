package structsql_test

import (
  "testing"

  "github.com/jimmc/jraceman/dbrepo/structsql"
)

func TestSelectSql(t *testing.T) {
  sql, targets := structsql.SelectSql("foo", foo)
  if got, want := sql,
      "SELECT id,num,required,optional,barid from foo"
      got != want {
    t.Errorf("SelectSql: got %v, want %v", got, want)
  }
  if got, want := len(targets), 5; got != want {
    t.Errorf("Wrong number of targets: got %d, want %d", got, want)
  }
}

func TestFindByIDSql(t *testing.T) {
  sql, targets := structsql.FindByIDSql("foo", foo)
  if got, want := sql, 
      "SELECT id,num,required,optional,barid from foo where id=?;";
      got != want {
    t.Errorf("FindByIDSql: got %v, want %v", got, want)
  }
  if got, want := len(targets), 5; got != want {
    t.Errorf("Wrong number of targets: got %d, want %d", got, want)
  }
}

func TestListSql(t *testing.T) {
  sql, _ := structsql.ListSql("foo", foo, 5, 10)
  if got, want := sql,
      "SELECT id,num,required,optional,barid from foo limit 10 offset 5"
      got != want {
    t.Errorf("ListSql: got %v, want %v", got, want)
  }
}
