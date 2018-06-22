package structsql

import (
  "database/sql"
  "fmt"
  "log"
  "strings"

  "github.com/jimmc/jracemango/dbrepo/strsql"
)

// Placeholder is ? for MySQL,$N for PostgreSQL,
// SQLite uses either of those, Oracle is :param1

// UpgradeTable upgrades a table from its current state to match
// what CreateTable would create.
func UpgradeTable(db *sql.DB, tableName string, entity interface{}, dryrun bool) (bool, string, error) {
  tableColumns, err := TableColumns(db, tableName)
  if err != nil {
    return false, "", fmt.Errorf("error getting columns for table %s: %v", tableName, err)
  }
  upgradeSql, err := UpgradeTableSql(tableName, entity, tableColumns)
  if err != nil {
    return false, "", err
  }
  if upgradeSql == "" {
    // Table is up to date.
    return true, "", nil
  }
  if dryrun {
    return false, upgradeSql, nil
  }
  _, err = db.Exec(upgradeSql)
  return false, upgradeSql, err
}

// UpgradeTableSql generates an SQL a create or upgrade command using
// the fields of the given struct. For each field in the struct:
//   * The field name is converted to lower case for the column name.
//   * int and string fields are declared as that same type column.
//   * The id field is declared as primary key.
//   * Non-pointer fields are declared as not null.
//   * Field names ending in ID are declared as foreign key references to the
//     id field of a table whose name matches the first part of the field name
// If the table does not exist, a CREATE TABLE statement is generated.
// If the table exists, an ALTER TABLE statement is generated to add any
// missing columns.
func UpgradeTableSql(tableName string, entity interface{}, tableColumns []ColumnInfo) (string, error) {
  // TODO - if the table does not exist, return the CREATE TABLE statement.
  columnInfos := ColumnInfos(entity)
  log.Printf("tableColumns for %s: %v", tableName, tableColumns)
  log.Printf("columnInfos for %s: %v", tableName, columnInfos)
  diffs := DiffColumnInfos(tableColumns, columnInfos)
  log.Printf("UpgradeTablesSql diffs for %s: %v", tableName, diffs)
  if len(diffs.Change) != 0 {
    // We don't know how to change columns, so this is an error.
    changedColNames := make([]string, len(diffs.Change))
    for i, cc := range diffs.Change {
      changedColNames[i] = cc[0].Name
    }
    return "", fmt.Errorf("Don't know how to change table %s columns %v", tableName, changedColNames)
  }
  if len(diffs.Add) == 0 {
    // No changes needed.
    return "", nil
  }
  // Create ALTER TABLE ADD COLUMN command for each column in diffs.Add.
  alterSql := ""
  for _, colInfo := range diffs.Add {
    columnSpec := colInfo.Name + " " + colInfo.Type
    if colInfo.Name == "id" {
      columnSpec = columnSpec + " primary key"
    } else {
      if colInfo.Required {
        columnSpec = columnSpec + " not null"
      }
      if colInfo.IsForeignKey {
        columnSpec = columnSpec + " references " + colInfo.FKTable + "(id)"
      }
    }
    alterColSql := "ALTER TABLE " + tableName + " ADD COLUMN " + columnSpec + "; "
    alterSql = alterSql + alterColSql;
  }
  log.Printf("UpgradeTableSql: %v\n", alterSql)
  return alterSql, nil
}

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
