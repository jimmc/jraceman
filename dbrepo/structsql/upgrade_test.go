package structsql_test

import (
  "strings"
  "testing"

  dbtest "github.com/jimmc/jracemango/dbrepo/test"
  "github.com/jimmc/jracemango/dbrepo/structsql"
)

func TestCreateOrUpdateTableSqlExists(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  tableSql, err := structsql.CreateOrUpgradeTableSql(db, "test", foo, []structsql.ColumnInfo{}, nil)
  if err != nil {
    t.Fatalf("Error generating sql: %v", err)
  }
  if !strings.HasPrefix(tableSql, "ALTER TABLE") {
    t.Errorf("Expected ALTER TABLE statement, got %v", tableSql)
  }
}

func TestCreateOrUpdateTableSqlNotExists(t *testing.T) {
  db, err := dbtest.EmptyDb()
  if err != nil {
    t.Fatal(err.Error())
  }
  tableSql, err := structsql.CreateOrUpgradeTableSql(db, "foo", foo, []structsql.ColumnInfo{}, nil)
  if err != nil {
    t.Fatalf("Error generating sql: %v", err)
  }
  if !strings.HasPrefix(tableSql, "CREATE TABLE") {
    t.Errorf("Expected CREATE TABLE statement, got %v", tableSql)
  }
}

func TestUpgradeTableSql(t *testing.T) {
  db, err := dbtest.DbWithSetupString(`
    CREATE TABLE foo(id string primary key, num int not null, required string not null)
  `)
  if err != nil {
    t.Fatal(err.Error())
  }
  expectedSql := `ALTER TABLE foo ADD COLUMN optional string; ALTER TABLE foo ADD COLUMN barid string references bar(id); `
  tableColumns, err := structsql.TableColumns(db, "foo")
  if err != nil {
    t.Fatal(err.Error())
  }
  tableSql, err := structsql.CreateOrUpgradeTableSql(db, "foo", foo, tableColumns, nil)
  if err != nil {
    t.Fatalf("Error generating sql: %v", err)
  }
  if got, want := tableSql, expectedSql; got != want {
    t.Errorf("UpgradeTable sql: got %v, want %v", got, want)
  }
}
