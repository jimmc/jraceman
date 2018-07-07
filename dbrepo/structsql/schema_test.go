package structsql_test

import (
  "reflect"
  "testing"

  dbtest "github.com/jimmc/jracemango/dbrepo/test"
  "github.com/jimmc/jracemango/dbrepo/structsql"
)

func TestColumnSpec(t *testing.T) {
  type entry struct {
    info structsql.ColumnInfo
    spec string
  }
  tests := []entry{
    {
      info: structsql.ColumnInfo{Name:"a", Type:"int", Required:false},
      spec: "a int",
    },
    {
      info: structsql.ColumnInfo{Name:"b", Type:"string", Required:true},
      spec: "b string not null",
    },
    {
      info: structsql.ColumnInfo{Name:"c", Type:"string", IsForeignKey:true, FKTable:"bar"},
      spec: "c string references bar(id)",
    },
  }
  for _, e := range tests {
    if got, want := structsql.ColumnSpec(e.info), e.spec; got != want {
      t.Errorf("ColumnSpec: got '%v', want '%v' for %v", got, want, e.info)
    }
  }
}

func TestTableExistsFalse(t *testing.T) {
  db, err := dbtest.EmptyDb()
  if err != nil {
    t.Fatal(err.Error())
  }
  exists, err := structsql.TableExists(db, "nosuchtable")
  if err != nil {
    t.Fatal(err.Error())
  }
  if exists {
    t.Errorf("Expected table not to exist")
  }
}

func TestTableExistsTrue(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  exists, err := structsql.TableExists(db, "test")
  if err != nil {
    t.Fatal(err.Error())
  }
  if !exists {
    t.Errorf("Expected test table to exist")
  }
}

func TestTableColumns(t *testing.T) {
  expectedColumns := []structsql.ColumnInfo{
    structsql.ColumnInfo{Name:"id", Type:"string", Required:true, IsForeignKey:false, FKTable:""},
    structsql.ColumnInfo{Name:"n", Type:"int", Required:false, IsForeignKey:false, FKTable:""},
    structsql.ColumnInfo{Name:"s", Type:"string", Required:false, IsForeignKey:false, FKTable:""},
    structsql.ColumnInfo{Name:"fooid", Type:"string", Required:false, IsForeignKey:true, FKTable:"foo"},
    structsql.ColumnInfo{Name:"amtpaid", Type:"int", Required:false, IsForeignKey:false, FKTable:""},
  }
  db, err := dbtest.DbWithSetupString(`
      CREATE table test(id string, n int, s string, fooid string references foo(id), amtpaid int);
      `)
  if err != nil {
    t.Fatal(err.Error())
  }
  columnInfos, err := structsql.TableColumns(db, "test")
  if err != nil {
    t.Fatal(err.Error())
  }
  if got, want := columnInfos, expectedColumns; !reflect.DeepEqual(got, want) {
    t.Errorf("ColumnInfo: got %v, want %v", got, want)
  }
}
