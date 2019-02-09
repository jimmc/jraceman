package ixport

import (
  "bufio"
  "fmt"
  "io"
  "log"
  "strconv"
  "strings"
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
  rowRepo RowRepo
  lineno int
  tableName string
  columnNames []string
  v1TableName string
  v1ColumnNames []string
  idIndex int
  counts ImporterCounts
  importVersion int
  v1v2map map[string]*tableMapValue
}

type RowRepo interface {
  Read(table string, columns []string, ID string) ([]interface{}, error)
  Insert(table string, columns[]string, values []interface{}, ID string) error
  Update(table string, columns[]string, values []interface{}, ID string) error
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

func (im *Importer) Import(reader io.Reader) error {
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
      log.Printf("Import: %d lines processed\n", im.lineno)
    }
    */
  }
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
  /* TODO - implement include, maybe sql, sqlexpect, sqlcheck
  case "include":
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
  im.tableName = tableName
  im.columnNames = []string{}
  im.idIndex = -1
  log.Printf("Import: table %s\n", im.tableName)
  return nil
}

func (im *Importer) setColumns(columns string) error {
  columnNames := strings.Split(columns, ",")
  hasID := false
  im.idIndex = -1
  for n, _ := range columnNames {
    columnNames[n] = strings.TrimPrefix(strings.TrimSuffix(columnNames[n], `"`), `"`)
    if columnNames[n] == "id" {
      hasID = true
      im.idIndex = n
    }
  }
  if !hasID {
    return fmt.Errorf("id column is required but missing on line %d", im.lineno)
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
    log.Printf("Import: translate v1 table %s to v2 table %s\n", im.v1TableName, im.tableName)
  }
  return nil
}

// ImportDataLine processes a line with field data values.
func (im *Importer) importDataLine(line string) error {
  if im.idIndex < 0 {
    return fmt.Errorf("id column has not been specified")
  }
  s := NewQuotedScanner(line)
  tokens, err := s.CommaSeparatedTokens()
  if err != nil {
    return err
  }
  values := s.TokensToValues(tokens)
  // TODO - look for strings of the form "{dt 'sss'}" (or d or t) as representing
  // timestamp (datetime), date or time.

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
  ID, ok := values[im.idIndex].(string)
  if !ok {
    return fmt.Errorf("id value must be a string, line %d", im.lineno)
  }
  // Look up the existing row to see if it exists and whether we are changing it.
  existingValues, err := im.rowRepo.Read(im.tableName, im.columnNames, ID)
  if err != nil {
    return fmt.Errorf("error retrieving existing data for %s[%s]: %v",
        im.tableName, values[im.idIndex], err)
  }
  if existingValues == nil {
    isNew = true
  } else {
    if len(values) != len(existingValues) {
      return fmt.Errorf("Wrong number of values to compare against existing: got %d, expected %d", len(values), len(existingValues))
    }
    // The row exists, look to see if any of our fields represent changes.
    for i := 0; i < len(values); i++ {
      if values[i] != existingValues[i] {
        diffColumns = append(diffColumns, im.columnNames[i])
        diffValues = append(diffValues, values[i])
        /*
        log.Printf("column:%s old:%v(%T) new:%v(%T)",
            im.columnNames[i], existingValues[i], existingValues[i], values[i], values[i])
        */
      }
    }
  }

  if isNew {
    if err := im.rowRepo.Insert(im.tableName, im.columnNames, values, ID); err != nil {
      return err
    }
    im.counts.inserted++
  } else if len(diffColumns) > 0 {
    if err := im.rowRepo.Update(im.tableName, diffColumns, diffValues, ID); err != nil {
      return err
    }
    im.counts.updated++
  } else {
    // No change to the existing data.
    im.counts.unchanged++
  }

  return nil
}
