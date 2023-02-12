package structsql

import (
  "github.com/jimmc/jraceman/dbrepo/compat"
)

// DeleteByID deletes a record by ID.
func DeleteByID(db compat.DBorTx, tableName, ID string) error {
  sql := DeleteByIDSql(tableName)
  res, err := db.Exec(sql, ID)
  return RequireOneResult(res, err, "Deleted", tableName, ID)
}

// DeleteByIDSql generates an SQL DELETE statement.
func DeleteByIDSql(tableName string) string {
  sql := "delete from " + tableName + " where id=?;"
  return sql
}
