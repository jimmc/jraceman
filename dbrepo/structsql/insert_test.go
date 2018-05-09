package structsql_test

import (
  "testing"

  "github.com/jimmc/jracemango/dbrepo/structsql"
)

func TestInsertSql(t *testing.T) {
  sql, values := structsql.InsertSql("foo", foo)
  if got, want := sql,
      "INSERT into foo(id,num,required,optional) values (?,?,?,?);";
      got != want {
    t.Errorf("InsertSql: got %v, want %v", got, want)
  }
  if got, want := len(values), 4; got != want {
    t.Errorf("Wrong number of values: got %d, want %d", got, want)
  }
}
