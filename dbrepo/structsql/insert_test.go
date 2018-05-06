package structsql

import (
  "testing"
)

func TestInsertSql(t *testing.T) {
  sql, values := InsertSql("foo", foo)
  if got, want := sql,
      "INSERT into foo(id,num,required,optional) values (?,?,?,?);";
      got != want {
    t.Errorf("InsertSql: got %v, want %v", got, want)
  }
  if got, want := len(values), 4; got != want {
    t.Errorf("Wrong number of values: got %d, want %d", got, want)
  }
}
