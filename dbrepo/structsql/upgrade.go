package structsql

import (
  "database/sql"
  "fmt"
  "log"
)

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
