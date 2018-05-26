package strsql

import (
  "database/sql"
)

// QueryResults provides generic results for an SQL query.
type QueryResults struct {
  Columns []*ColumnInfo
  Rows [][]interface{}
}

// ColumnInfo provides information about one of the columns in the results of a query.
type ColumnInfo struct {
  Name string
  Type string
}

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

// QueryStarAndCollect issues a Query for the given arbitrary sql and returns the results.
// It is intended for cases where the type and column count of the result is unknown,
// such as "SELECT * from sometable".
func QueryStarAndCollect(db *sql.DB, sql string, queryValues ...interface{}) (*QueryResults, error) {
  rows, err := db.Query(sql, queryValues...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  columnTypes, err := rows.ColumnTypes()
  if err != nil {
    return nil, err
  }
  columnInfo := make([]*ColumnInfo, len(columnTypes))
  for c, col := range columnTypes {
    dbType := col.DatabaseTypeName()
    if isStringType(dbType) {
      dbType = "string"
    }
    columnInfo[c] = &ColumnInfo{
      Name: col.Name(),
      Type: dbType,
    }
  }
  data := make([][]interface{}, 0)
  rowTargets := make([]interface{}, len(columnTypes))
    // Each field of rowTargets gets pointed to each field of rowData
    // before each row scan.
  for rows.Next() {
    rowData := make([]interface{}, len(columnTypes))
    for c := range rowTargets {
      rowTargets[c] = &rowData[c]
    }
    err := rows.Scan(rowTargets...)
    if err != nil {
      return nil, err
    }
    for c, col := range columnTypes {
      colDbTypeName := col.DatabaseTypeName()
      if isStringType(colDbTypeName) {
        fieldBytes, ok := rowData[c].([]byte)
        if ok {
          rowData[c] = string(fieldBytes)
        }
      }
    }
    data = append(data, rowData)
  }
  if err := rows.Err(); err != nil {
    return nil, err
  }
  results := &QueryResults{
    Columns: columnInfo,
    Rows: data,
  }
  return results, nil
}

func isStringType(ctype string) bool {
  stringTypes := []string{"string", "text", "varchar"}
  for _, t := range stringTypes {
    if ctype == t {
      return true
    }
  }
  return false
}
