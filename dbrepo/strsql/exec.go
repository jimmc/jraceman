package strsql

import (
  "regexp"
  "strings"

  "github.com/jimmc/jraceman/dbrepo/compat"
)

// ExecMulti executes multiple sql statements from a string.
// It strips out all carriage returns, breaks the string into segments at double-newlines, removes lines
// starting with a "#" (as comments), and separately executes
// each segment. If any segment returns an error, it stop executing
// and returns that error.
func ExecMulti(db compat.DBorTx, sql string) error {
  re := regexp.MustCompile("\r")
  sql = re.ReplaceAllString(sql, "")
  segments := strings.Split(sql, "\n\n")
  for _, segment := range segments {
    if err := ExecSegment(db, segment); err != nil {
      return err
    }
  }
  return nil
}

// ExecSegment executes a single sql statement from a string.
// It removes lines starting with "#" (as comments).
func ExecSegment(db compat.DBorTx, segment string) error {
  lines := strings.Split(segment, "\n")
  sqlLines := make([]string, 0)
  for _, line := range lines {
    if !strings.HasPrefix(line, "#") {
      sqlLines = append(sqlLines, line)
    }
  }
  segment = strings.Join(sqlLines, "\n")
  _, err := db.Exec(segment)
  return err
}
