package structsql

import (
  "fmt"

  "github.com/jimmc/jraceman/dbrepo/conn"

  "github.com/golang/glog"
)

// UpgradeTable upgrades a table from its current state to match
// what CreateTable would create. If the table does not exist,
// it creates the table.
func UpgradeTable(db conn.DB, tableName string, entity interface{}, dryrun bool) (bool, string, error) {
  return UpgradeTableWithUpdater(db, tableName, entity, dryrun, nil)
}

// UpgradeTableWithUpdater upgrades a table from its current state to match
// what CreateTableWithUpdater would create. If the table does not exist,
// it creates the table.
func UpgradeTableWithUpdater(db conn.DB, tableName string, entity interface{}, dryrun bool, updater ColumnInfosUpdater) (bool, string, error) {
  tableColumns, err := TableColumns(db, tableName)
  if err != nil {
    return false, "", fmt.Errorf("error getting columns for table %s: %v", tableName, err)
  }
  tableSql, err := CreateOrUpgradeTableSql(db, tableName, entity, tableColumns, updater)
  if err != nil {
    return false, "", err
  }
  if tableSql == "" {
    // Table is up to date.
    return true, "", nil
  }
  if dryrun {
    return false, tableSql, nil
  }
  _, err = db.Exec(tableSql)
  return false, tableSql, err
}

// CreateOrUpgradeTableSql checks to see whether the table already exists,
// and returns either a CREATE TABLE statement if it does not exist, or
// the ALTER TABLE statements for column changes if it does exist.
func CreateOrUpgradeTableSql(db conn.DB, tableName string, entity interface{}, tableColumns []ColumnInfo, updater ColumnInfosUpdater) (string, error) {
  exists, err := TableExists(db, tableName)
  if err != nil {
    return "", err
  }
  columnInfos := ColumnInfos(entity)
  if updater != nil {
    glog.V(2).Infof("CreateTableSql calling UpdateColumnInfo for table %s", tableName)
    tableColumns = updater.UpdateColumnInfos(tableColumns)
    columnInfos = updater.UpdateColumnInfos(columnInfos)
  }
  if exists {
    return UpgradeTableSql(tableName, columnInfos, tableColumns)
  } else {
    return CreateTableSqlFromColumnInfos(tableName, columnInfos), nil
  }
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
func UpgradeTableSql(tableName string, columnInfos, tableColumns []ColumnInfo) (string, error) {
  glog.V(1).Infof("tableColumns for %s: %v", tableName, tableColumns)
  glog.V(1).Infof("columnInfos for %s: %v", tableName, columnInfos)
  diffs := DiffColumnInfos(tableColumns, columnInfos)
  glog.V(1).Infof("UpgradeTablesSql diffs for %s: %v", tableName, diffs)
  if len(diffs.Change) != 0 {
    // We don't know how to change columns, so this is an error.
    changedColNames := make([]string, len(diffs.Change))
    for i, cc := range diffs.Change {
      changedColNames[i] = cc[0].Name
    }
    return "", fmt.Errorf("Don't know how to change table %s columns %v; diffs %#v", tableName, changedColNames, diffs.Change)
  }
  if len(diffs.Add) == 0 {
    // No changes needed.
    return "", nil
  }
  // Create ALTER TABLE ADD COLUMN command for each column in diffs.Add.
  alterSql := ""
  for _, colInfo := range diffs.Add {
    columnSpec := ColumnSpec(colInfo)
    alterColSql := "ALTER TABLE " + tableName + " ADD COLUMN " + columnSpec + "; "
    alterSql = alterSql + alterColSql;
  }
  glog.V(1).Infof("UpgradeTableSql: %v\n", alterSql)
  return alterSql, nil
}
