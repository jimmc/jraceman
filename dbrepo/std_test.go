package dbrepo

import (
  "errors"
  "reflect"
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
  if got, want := stdCreateTableSqlFromStruct("foo", foo),
      "CREATE TABLE foo(id string primary key, num int not null, " +
      "required string not null, optional string, opt2 string);";
      got != want {
    t.Errorf("CreateTableSql: got %v, want %v", got, want)
  }
}

func TestSelectSql(t *testing.T) {
  sql, targets := stdSelectSqlFromStruct("foo", foo)
  if got, want := sql,
      "SELECT id,num,required,optional,opt2 from foo"
      got != want {
    t.Errorf("SelectSql: got %v, want %v", got, want)
  }
  if got, want := len(targets), 5; got != want {
    t.Errorf("Wrong number of targets: got %d, want %d", got, want)
  }
}

func TestFindByIDSql(t *testing.T) {
  sql, targets := stdFindByIDSqlFromStruct("foo", foo)
  if got, want := sql, 
      "SELECT id,num,required,optional,opt2 from foo where id=?;";
      got != want {
    t.Errorf("FindByIDSql: got %v, want %v", got, want)
  }
  if got, want := len(targets), 5; got != want {
    t.Errorf("Wrong number of targets: got %d, want %d", got, want)
  }
}

func TestInsertSqlFromStruct(t *testing.T) {
  sql, values := stdInsertSqlFromStruct("foo", foo)
  if got, want := sql,
      "INSERT into foo(id,num,required,optional) values (?,?,?,?);";
      got != want {
    t.Errorf("InsertSql: got %v, want %v", got, want)
  }
  if got, want := len(values), 4; got != want {
    t.Errorf("Wrong number of values: got %d, want %d", got, want)
  }
}

func TestDeleteByIDSql(t *testing.T) {
  if got, want := stdDeleteByIDSql("foo"), "delete from foo where id=?;"; got != want {
    t.Errorf("DeleteByIDSql: got %v, want %v", got, want)
  }
}

func TestModsToSql(t *testing.T) {
  mods := map[string]interface{}{
    "Num": 456,
    "Required": "qqq",
    "Opt2": nil,
  }
  sql, values := modsToSql("foo", mods, "123")
  if got, want := sql, "update foo set num = ?, opt2 = NULL, required = ? where id = ?;"; got != want {
    t.Errorf("Update sql: got %v, want %v", got, want)
  }
  // Values should not include Opt2, but does include ID
  if got, want := len(values), 3; got != want {
    t.Errorf("Number of values: got %d, want %d", got, want)
  }
  expectedValues := []interface{}{456, "qqq", "123"}
  if got, want := values, expectedValues; !reflect.DeepEqual(got, want) {
    t.Errorf("Values: got %v, want %v", got, want)
  }
}

type oneResultTester struct {
  num int64
  err error
}
func (o *oneResultTester) RowsAffected() (int64, error) {
  return o.num, o.err
}

func TestRequireOneResult(t *testing.T) {
  oZero := &oneResultTester{0, nil}
  oOne := &oneResultTester{1, nil}
  oTwo := &oneResultTester{2, nil}
  oErr := &oneResultTester{0, errors.New("Test error")}

  if got, want := requireOneResult(oOne, nil, "Tested", "foo", "123"), error(nil); got != want {
    t.Errorf("Happy path: got %v, want %v", got, want)
  }
  if got, want := requireOneResult(oOne, errors.New("Test error"), "Tested", "foo", "123"), errors.New("Test error"); got.Error() != want.Error() {
    t.Errorf("With passed-in error: got %v, want %v", got, want)
  }
  if got, want := requireOneResult(oErr, nil, "Tested", "foo", "123"), errors.New("Test error"); got.Error() != want.Error() {
    t.Errorf("With sql error: got %v, want %v", got, want)
  }
  if got, want := requireOneResult(oZero, nil, "Tested", "foo", "123"), errors.New("Wrong-count error"); got == nil {
    t.Errorf("With count==0: got %v, want %v", got, want)
  }
  if got, want := requireOneResult(oTwo, nil, "Tested", "foo", "123"), errors.New("Wrong-count error"); got == nil {
    t.Errorf("With count==0: got %v, want %v", got, want)
  }
}
