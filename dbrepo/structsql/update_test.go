package structsql

import (
  "reflect"
  "testing"
)

func TestModsToSql(t *testing.T) {
  mods := map[string]interface{}{
    "Num": 456,
    "Required": "qqq",
    "Opt2": nil,
  }
  sql, values := ModsToSql("foo", mods, "123")
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
