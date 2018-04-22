package dbrepo

import (
  "database/sql"
  "fmt"
)

// RequireOneResult gets the result of sql.Stmt.Exec and verifies that it
// affected exactly one row, which should be the case for operations that
// use the entity ID. The action string is used in error messages, and
// should be capitalized and past tense, such as "Deleted".
// The entityType should the name of the entity, such as "site".
func requireOneResult(res sql.Result, err error, action, entityType, ID string) error {
  if err != nil {
    return err
  }
  rowCnt, err := res.RowsAffected()
  if err != nil {
    return err
  }
  if rowCnt == 0 {
    return fmt.Errorf("%s no %s rows for ID %s", action, entityType, ID)
  }
  if rowCnt > 1 {
    return fmt.Errorf("%s %d %s rows for ID %s", action, rowCnt, entityType, ID)
  }
  return nil
}
