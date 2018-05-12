package structsql_test

import (
  "reflect"
  "testing"

  "github.com/jimmc/jracemango/dbrepo/dbtest"
  "github.com/jimmc/jracemango/dbrepo/structsql"
)

type Foo struct {
  ID string
  Num int
  Required string
  Optional *string
  BarID *string
}

var str = "opt"
var foo = &Foo{
  ID: "X1",
  Num: 123,
  Required: "rrr",
  Optional: &str,
}

func TestCreateTableSql(t *testing.T) {
  if got, want := structsql.CreateTableSql("foo", foo),
      "CREATE TABLE foo(id string primary key, num int not null, " +
      "required string not null, optional string, barid string references bar(id));";
      got != want {
    t.Errorf("CreateTableSql: got %v, want %v", got, want)
  }
}

func TestColumnNames(t *testing.T) {
  columnNames := structsql.ColumnNames(foo)
  expectedNames:= []string{"id", "num", "required", "optional", "barid"}
  if got, want := columnNames, expectedNames; !reflect.DeepEqual(got, want) {
    t.Errorf("Column names: got %v, want %v", got, want)
  }
}

func TestCreateTable(t *testing.T) {
  db, err := dbtest.EmptyDb()
  if err != nil {
    t.Fatal(err.Error())
  }
  if err := structsql.CreateTable(db, "foo", &Foo{}); err != nil {
    t.Fatalf("Creating table: %v", err)
  }
}
