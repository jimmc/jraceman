package structsql_test

import (
  "testing"

  "github.com/jimmc/jracemango/dbrepo/dbtest"
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

func TestInsert(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  testRecord := dbtest.TestRecord{
    ID: "T25",
    N: 123,
    S: "xyz",
  }
  if err = structsql.Insert(db, "test", &testRecord, testRecord.ID); err != nil {
    t.Fatalf("Inserting: %v", err)
  }
  readRecord := dbtest.TestRecord{}
  if err = db.QueryRow(
      "SELECT n, s from test where id=?;", testRecord.ID).Scan(
      &readRecord.N, &readRecord.S); err != nil {
    t.Errorf("Reading: %f", err)
  }
  if got, want := readRecord.N, testRecord.N; got != want {
    t.Errorf("Field N: got %v, want %v", got, want)
  }
  if got, want := readRecord.S, testRecord.S; got != want {
    t.Errorf("Field S: got %v, want %v", got, want)
  }
}
