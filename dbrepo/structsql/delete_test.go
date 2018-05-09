package structsql_test

import (
  "testing"

  "github.com/jimmc/jracemango/dbrepo/structsql"
)

func TestDeleteByIDSql(t *testing.T) {
  if got, want := structsql.DeleteByIDSql("foo"), "delete from foo where id=?;"; got != want {
    t.Errorf("DeleteByIDSql: got %v, want %v", got, want)
  }
}
