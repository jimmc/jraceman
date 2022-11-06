package strsql

import (
  "database/sql"
  "fmt"
  "strings"
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

// QueryAndCollect issues a Query for the given sql, then iterates through
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
    dbType := strings.ToLower(col.DatabaseTypeName())
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
      colDbTypeName := strings.ToLower(col.DatabaseTypeName())
      // sqlite returns blank type names on pragma commands,
      // so try to convert those to strings when that works.
      if isStringType(colDbTypeName) || colDbTypeName == "" {
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

// QueryInt issues a Query for the given sql which is expected to have a single
// integer as a return value, such as "SELECT count(*) FROM ... WHERE ..."
func QueryInt(db *sql.DB, sql string, queryValues ...interface{}) (int, error) {
  rows, err := db.Query(sql, queryValues...)
  if err != nil {
    return 0, err
  }
  defer rows.Close()
  var n int
  if rows.Next() {
    err := rows.Scan(&n)
    if err != nil {
      return 0, err
    }
  } else {
    return 0, fmt.Errorf("Expected query to return a row")
  }
  if err = rows.Err(); err != nil {
    return 0, err
  }
  return n, nil
}

// QueryString issues a Query for the given sql which is expected to have a single
// string from a single record as a return value.
func QueryString(db *sql.DB, sql string, queryValues ...interface{}) (string, error) {
  rows, err := db.Query(sql, queryValues...)
  if err != nil {
    return "", err
  }
  defer rows.Close()
  var s string
  if rows.Next() {
    err := rows.Scan(&s)
    if err != nil {
      return "", err
    }
  } else {
    return "", fmt.Errorf("Expected query to return a row")
  }
  if err = rows.Err(); err != nil {
    return "", err
  }
  return s, nil
}

// QueryStrings issues a Query for the given sql and collects the first column of the
// result, which it assumes is a string.
func QueryStrings(db *sql.DB, sql string, queryValues ...interface{}) ([]string, error) {
  rows, err := db.Query(sql, queryValues...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  ss := make([]string, 0)
  var s string
  for rows.Next() {
    err := rows.Scan(&s)
    if err != nil {
      return nil, err
    }
    ss = append(ss, s)
  }
  if err = rows.Err(); err != nil {
    return nil, err
  }
  return ss, nil
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
