package structsql

import (
  "database/sql"
  "fmt"
  "log"
  "strings"

  "github.com/jimmc/jracemango/dbrepo/strsql"
)

// TableColumns collects details about the current columns in the database
// for the specified table, and returns them in the same form as ColumnInfos.
func TableColumns(db *sql.DB, tableName string) ([]ColumnInfo, error) {
  // TODO - this is specific to SQLite
  colSql := `select name, type, "notnull", pk from pragma_table_info(?)`
  results, err := strsql.QueryStarAndCollect(db, colSql, tableName)
  if err != nil {
    return nil, fmt.Errorf("error getting column info for table %s: %v", tableName, err)
  }
  log.Printf("colinfo query results table %s: %v", tableName, results)
  // The SQLite pragma command seems not to return a data type for the pragma columns,
  // so we manually do some conversions for the columns we know are strings.
  for r, row := range results.Rows {
    for c, col := range results.Columns {
      if col.Name == "name" || col.Name == "type" {
        fieldBytes, ok := row[c].([]byte)
        if ok {
          results.Rows[r][c] = string(fieldBytes)
        }
      }
    }
  }
  log.Printf("colinfo query results stringified: %v", results)
  // Now convert from a query result to []ColumnInfo
  colInfos := make([]ColumnInfo, len(results.Rows))
  for r, row := range results.Rows {
    name := row[0].(string)
    req := false;
    if row[2].(int64) != 0  || name == "id" {
      req = true;
    }
    // We don't have the info about whether it is a foreign key,
    // so just use the same rules as when we parse an entity.
    fk := false
    fkReference := ""
    if name != "id" && strings.HasSuffix(name, "id") {
      fk = true
      fkReference = strings.TrimSuffix(name, "id")
    }
    colType := row[1].(string)
    if colType == "boolean" {
      colType = "bool"
    }
    colInfos[r] = ColumnInfo{
      Name: name,
      Type: colType,
      Required: req,
      IsForeignKey: fk,
      FKTable: fkReference,
    }
  }
  log.Printf("colinfos returned: %v", colInfos)
  return colInfos, nil
}
