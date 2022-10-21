package structsql

import (
  "database/sql"
  "fmt"
  "strings"

  "github.com/jimmc/jraceman/dbrepo/strsql"
)

// ColumnSpec generates an sql column specification, suitable for use
// in a CREATE TABLE or ALTER TABLE command, for the given ColumnInfo.
func ColumnSpec(colInfo ColumnInfo) string {
  columnSpec := colInfo.Name + " " + colInfo.Type
  if colInfo.Name == "id" {
    columnSpec = columnSpec + " primary key"
  } else {
    if colInfo.Required {
      columnSpec = columnSpec + " not null"
    }
    if colInfo.HasDefault {
      columnSpec = columnSpec + " default " + colInfo.DefaultAsString
    }
    if colInfo.IsForeignKey {
      columnSpec = columnSpec + " references " + colInfo.FKTable + "(id)"
    }
  }
  return columnSpec
}

// TableExists returns true of the specified table exists in the database.
func TableExists(db *sql.DB, tableName string) (bool, error) {
  // TODO - this is specific to SQLite
  colSql := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
  results, err := strsql.QueryStarAndCollect(db, colSql, tableName)
  if err != nil {
    return false, fmt.Errorf("error getting info for table %s: %v", tableName, err)
  }
  for _, _ = range results.Rows {
    return true, nil    // If we have any results, the table exists.
  }
  return false, nil
}

// TableColumns collects details about the current columns in the database
// for the specified table, and returns them in the same form as ColumnInfos.
func TableColumns(db *sql.DB, tableName string) ([]ColumnInfo, error) {
  // TODO - this is specific to SQLite
  colSql := `select name, lower(type), "notnull", pk from pragma_table_info(?)`
  results, err := strsql.QueryStarAndCollect(db, colSql, tableName)
  if err != nil {
    return nil, fmt.Errorf("error getting column info for table %s: %v", tableName, err)
  }
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
    // TODO - might want to switch this to using the table definition that
    // includes the foreign key stuff so we don't need this hack.
    if name != "id" && strings.HasSuffix(name, "id") && !strings.HasSuffix(name, "paid") {
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
  return colInfos, nil
}
