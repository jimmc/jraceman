package structsql_test

import (
  "reflect"
  "testing"

  "github.com/jimmc/jracemango/dbrepo/dbtest"
  "github.com/jimmc/jracemango/dbrepo/structsql"
)

func TestModsToSql(t *testing.T) {
  mods := map[string]interface{}{
    ".Num": 456,                // Make sure leading dot is stripped off.
    "Required": "qqq",
    "Opt2": nil,
  }
  sql, values := structsql.ModsToSql("foo", mods, "123")
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

type testUpdateEntity struct {
  ID string
  N int
  S string
}

func TestUpdateById(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()
  ID := "T2"
  mods := map[string]interface{}{
    "n": 123,
  }
  if err = structsql.UpdateByID(db, "test", mods, ID); err != nil {
    t.Errorf("Error updating: %v", err)
  }
  testEnt := &testUpdateEntity{}
  sql, targets := structsql.FindByIDSql("test", testEnt)
  if err := db.QueryRow(sql, ID).Scan(targets...); err != nil {
    t.Errorf("Scanning row: %v", err)
  }
  if got, want := testEnt.N, 123; got != want {
    t.Errorf("Wrong value after update: got %v, want %v", got, want)
  }
}
