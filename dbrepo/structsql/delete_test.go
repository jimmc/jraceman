package structsql

import (
  "testing"
)

func TestDeleteByIDSql(t *testing.T) {
  if got, want := DeleteByIDSql("foo"), "delete from foo where id=?;"; got != want {
    t.Errorf("DeleteByIDSql: got %v, want %v", got, want)
  }
}
