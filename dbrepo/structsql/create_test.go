package structsql

import (
  "testing"
)

type Foo struct {
  ID string
  Num int
  Required string
  Optional *string
  Opt2 *string
}

var str = "opt"
var foo = &Foo{
  ID: "X1",
  Num: 123,
  Required: "rrr",
  Optional: &str,
}

func TestCreateTableSql(t *testing.T) {
  if got, want := CreateTableSql("foo", foo),
      "CREATE TABLE foo(id string primary key, num int not null, " +
      "required string not null, optional string, opt2 string);";
      got != want {
    t.Errorf("CreateTableSql: got %v, want %v", got, want)
  }
}
