package structsql_test

import (
  "strings"
  "testing"

  dbtest "github.com/jimmc/jraceman/dbrepo/test"
  "github.com/jimmc/jraceman/dbrepo/structsql"
)

func TestDeleteByIDSql(t *testing.T) {
  if got, want := structsql.DeleteByIDSql("foo"), "delete from foo where id=?;"; got != want {
    t.Errorf("DeleteByIDSql: got %v, want %v", got, want)
  }
}

func TestDeleteByID(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatalf(err.Error())
  }
  readRecord := dbtest.TestRecord{}
  if err := db.QueryRow(
      "SELECT n, s from test where id=?;", "T1").Scan(
      &readRecord.N, &readRecord.S); err != nil {
    t.Fatalf("Reading: %f", err)
  }
  if err := structsql.DeleteByID(db, "test", "T1"); err != nil {
    t.Fatalf("Deleting: %v", err)
  }
  err = db.QueryRow(
      "SELECT n, s from test where id=?;", "T1").Scan(
      &readRecord.N, &readRecord.S)
  if err == nil {
    t.Fatalf("Expected error trying to read row after delete")
  }
  if !strings.Contains(err.Error(), "no rows") {
    t.Errorf("Wrong error after delete, got %v, want 'no rows'", err)
  }
}
