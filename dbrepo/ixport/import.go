package ixport

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "path/filepath"
  "strconv"
  "strings"

  "github.com/golang/glog"
)

type ImporterCounts struct {
  inserted int
  updated int
  unchanged int
}

func (i ImporterCounts) Inserted() int {
  return i.inserted
}

func (i ImporterCounts) Updated() int {
  return i.updated
}

func (i ImporterCounts) Unchanged() int {
  return i.unchanged
}

type Importer struct {
  fileName string       // The current source filename.
  rowRepo RowRepo
  lineno int
  tableName string
  columnNames []string
  v1TableName string
  v1ColumnNames []string
  idIndex int
  nameIndex int
  counts ImporterCounts
  tableStartCounts ImporterCounts
  importVersion int
  v1v2map map[string]*tableMapValue
}

type RowRepo interface {
  ReadByKey(table string, columns []string, keyName, key string) ([]interface{}, error)
  Insert(table string, columns[]string, values []interface{}, key string) error
  UpdateByKey(table string, columns[]string, values []interface{}, keyName, key string) error
}

func NewImporter(rowRepo RowRepo) *Importer {
  im := &Importer{
    rowRepo: rowRepo,
  }
  im.initTableMaps()
  return im
}

func (im *Importer) TableName() string {
  return im.tableName
}

func (im *Importer) ColumnNames() []string {
  return im.columnNames
}

func (im *Importer) Counts() ImporterCounts {
  return im.counts
}

func (im *Importer) ImportFile(importFile string) error {
  oldFileName := im.fileName
  defer func(){ im.fileName = oldFileName }()
  im.fileName = importFile
  inFile, err := os.Open(importFile)
  if err != nil {
    return fmt.Errorf("error opening import input file %s: %v", importFile, err)
  }
  defer inFile.Close()

  glog.Infof("Importing from %s\n", importFile)
  err = im.importReader(inFile)
  glog.Infof("Done importing from %s\n", importFile)
  return err
}

func (im *Importer) importReader(reader io.Reader) error {
  im.reset()
  scanner := bufio.NewScanner(reader)
  for scanner.Scan() {
    im.lineno++         // First line is line 1.
    line := scanner.Text()
    if err := im.ImportLine(line); err != nil {
      return fmt.Errorf("line %d: %v", im.lineno, err)
    }
    /*
    if im.lineno % 100 == 0 {
      glog.Infof("Import: %d lines processed\n", im.lineno)
    }
    */
  }
  im.printPreviousTableStats()
  if err := scanner.Err(); err != nil {
    return err
  }
  return nil
}

func (im *Importer) reset() {
  im.lineno = 0
  im.tableName = ""
}

func (im *Importer) ImportLine(line string) error {
  line = strings.TrimSpace(line)
  if line == "" || strings.HasPrefix(line, "#") {
    return nil          // Blank line or comment line.
  }
  if strings.HasPrefix(line, "!") {
    return im.importModeLine(line)
  } else {
    return im.importDataLine(line)
  }
}

// ImportModeLine processes a line starting with "!".
func (im *Importer) importModeLine(line string) error {
  line = strings.TrimPrefix(line, "!")
  words := strings.Fields(line)
  if len(words) == 0 {
    return nil          // Ignore blank commands
  }
  if len(words) < 2 {
    return fmt.Errorf("argument is required for command %s", words[0])
  }
  switch words[0] {
  case "exportVersion":
    return im.setExportVersion(words[1])
  case "appInfo":
    return im.setAppInfo(words[1])
  case "type":
    return im.setType(words[1])
  case "table":
    return im.setTable(words[1])
  case "columns":
    return im.setColumns(words[1])
  case "include":
    return im.importInclude(words[1])
  /* TODO - implement maybe sql, sqlexpect, sqlcheck
  case "sql":
  case "sqlexpect":
  case "sqlcheck":
  */
  default:
    return fmt.Errorf("unknown command %s", words[0])
  }
}

func (im *Importer) setExportVersion(version string) error {
  verNum, err := strconv.Atoi(version)
  if err != nil {
    return fmt.Errorf("Bad version format: %v", err)
  }
  if verNum < 1 || verNum > 2 {
    return fmt.Errorf("Unknown version number %d, must be 1 or 2", verNum)
  }
  im.importVersion = verNum;
  return nil
}

func (im *Importer) setAppInfo(appInfo string) error {
  // TODO - what should we do with the appInfo?
  return nil
}

func (im *Importer) setType(appType string) error {
  // TODO - what should we do with the appType?
  return nil
}

func (im *Importer) setTable(tableName string) error {
  im.printPreviousTableStats()
  im.tableStartCounts = im.counts
  im.tableName = tableName
  im.columnNames = []string{}
  im.idIndex = -1
  im.nameIndex = -1
  glog.Infof("Import: table %s\n", im.tableName)
  return nil
}

func (im *Importer) printPreviousTableStats() {
  if im.tableName == "" {
    return              // No previous table.
  }
  tableInserted := im.counts.inserted - im.tableStartCounts.inserted
  tableUpdated := im.counts.updated - im.tableStartCounts.updated
  tableUnchanged := im.counts.unchanged - im.tableStartCounts.unchanged
  glog.Infof("Import: for table %s: inserted %d, updated %d, unchanged %d",
      im.tableName, tableInserted, tableUpdated, tableUnchanged)
}

func (im *Importer) setColumns(columns string) error {
  columnNames := strings.Split(columns, ",")
  hasID := false
  hasName := false
  im.idIndex = -1
  im.nameIndex = -1
  for n, _ := range columnNames {
    columnNames[n] = strings.TrimPrefix(strings.TrimSuffix(columnNames[n], `"`), `"`)
    if columnNames[n] == "id" {
      hasID = true
      im.idIndex = n
    }
    if columnNames[n] == "name" {
      hasName = true
      im.nameIndex = n
    }
  }
  if !hasID && !hasName {
    return fmt.Errorf("id or name column is required but missing on line %d", im.lineno)
  }
  im.columnNames = columnNames
  if im.importVersion == 1 {
    v2TableName, v2ColumnNames, err := im.translateNames1To2(im.tableName, im.columnNames)
    if err != nil {
      return err
    }
    im.v1TableName = im.tableName
    im.v1ColumnNames = im.columnNames
    im.tableName = v2TableName
    im.columnNames = v2ColumnNames
    glog.Infof("Import: translate v1 table %s to v2 table %s\n", im.v1TableName, im.tableName)
  }
  return nil
}

// ImportDataLine processes a line with field data values.
func (im *Importer) importDataLine(line string) error {
  if im.idIndex < 0 && im.nameIndex < 0 {
    return fmt.Errorf("id or name column has not been specified")
  }
  s := NewQuotedScanner(line)
  tokens, err := s.CommaSeparatedTokens()
  if err != nil {
    return err
  }
  values := s.TokensToValues(tokens)

  if im.importVersion == 1 {
    if len(values) != len(im.v1ColumnNames) {
      return fmt.Errorf("wrong number of fields on line %d, got %d for v1 column count %d",
          im.lineno, len(values), len(im.v1ColumnNames))
    }
  } else {
    if len(values) != len(im.columnNames) {
      return fmt.Errorf("wrong number of fields on line %d, got %d for column count %d",
          im.lineno, len(values), len(im.columnNames))
    }
  }

  if im.importVersion == 1 {
    values = im.translateValues1To2(values)
  }

  isNew := false
  diffColumns := []string{}
  diffValues := []interface{}{}
  keyIndex := im.idIndex
  keyName := "id"
  if keyIndex < 0 {
    keyIndex = im.nameIndex
    keyName = "name"
  }
  key, ok := values[keyIndex].(string)
  if !ok {
    return fmt.Errorf("id or name value must be a string, line %d", im.lineno)
  }
  // Look up the existing row to see if it exists and whether we are changing it.
  existingValues, err := im.rowRepo.ReadByKey(im.tableName, im.columnNames, keyName, key)
  if err != nil {
    return fmt.Errorf("error retrieving existing data for %s[%s]: %v",
        im.tableName, key, err)
  }
  if existingValues == nil {
    isNew = true
  } else {
    if len(values) != len(existingValues) {
      return fmt.Errorf("Wrong number of values to compare against existing: got %d, expected %d", len(values), len(existingValues))
    }
    // The row exists, look to see if any of our fields represent changes.
    for i := 0; i < len(values); i++ {
      if fieldChanged(values[i], existingValues[i]) {
        diffColumns = append(diffColumns, im.columnNames[i])
        diffValues = append(diffValues, values[i])
        /*
        glog.V(1).Infof("column:%s old:%v(%T) new:%v(%T)",
            im.columnNames[i], existingValues[i], existingValues[i], values[i], values[i])
        */
      }
    }
  }

  if isNew {
    if err := im.rowRepo.Insert(im.tableName, im.columnNames, values, key); err != nil {
      return err
    }
    im.counts.inserted++
  } else if len(diffColumns) > 0 {
    if err := im.rowRepo.UpdateByKey(im.tableName, diffColumns, diffValues, keyName, key); err != nil {
      return err
    }
    im.counts.updated++
  } else {
    // No change to the existing data.
    im.counts.unchanged++
  }

  return nil
}

// ImportInclude processes the !include directive to import a nested file.
func (im *Importer) importInclude(fileName string) error {
  if !filepath.IsAbs(fileName) {
    // If the included fileName is not absolute, make it relative to the current file.
    fileName = filepath.Join(filepath.Dir(im.fileName), fileName)
  }
  // Restore the line number in the current file after importing the included file.
  currentLineNo := im.lineno
  defer func(){ im.lineno = currentLineNo }()
  return im.ImportFile(fileName)
}

func fieldChanged(newVal interface{}, oldVal interface{}) bool {
  if newVal == oldVal {
    return false
  }
  // Look for values that compare differently, but we consider unchanged.
  switch newV := newVal.(type) {
  case float32:
    switch oldV := oldVal.(type) {
    case float64:
      return float64(newV) != oldV
    case int:
      return newV != float32(oldV)
    }
  case bool:
    switch oldV := oldVal.(type) {
    case int:
      return newV != (oldV != 0)
    }
  case string:
    switch oldV := oldVal.(type) {
    case int:
      return newV != strconv.Itoa(oldV)
    }
  }
  // Didn't find any compatible types.
  return true
}
