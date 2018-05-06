package structsql

import (
  "database/sql"
)

// QueryAndCollect issues a Query for the given sql, then interates through
// the returned rows. For each row, it retrieves the data into targets, then
// calls the collect function. The assumption is that the targets store the
// results into data that is accessible to the collect function.
func QueryAndCollect(db *sql.DB, sql string, targets []interface{}, collect func()) error {
  rows, err := db.Query(sql)
  if err != nil {
    return err
  }
  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(targets...)
    if err != nil {
      return err
    }
    collect()
  }
  return rows.Err()
}
