package dbrepo

import (
  "fmt"
  "io"
  "strings"

  "github.com/jimmc/jracemango/dbrepo/structsql"
)

func (r *Repos) Export(w io.Writer) error {
  if err := r.exportHeader(w); err != nil {
    return err
  }

  // The order of output of the tables is important: tables with
  // foreign keys should be after the tables the point to.

  if err := r.dbSite.Export(r, w); err != nil {
    return err
  }

  if err := r.dbArea.Export(r, w); err != nil {
    return err
  }

  return nil
}

func (r *Repos) exportHeader(w io.Writer) error {
  io.WriteString(w, "#!jraceman -import\n")
  io.WriteString(w, "!exportVersion 2\n")
  io.WriteString(w, "!appInfo JRaceman v2.0.0\n")       // TODO - get real version
  io.WriteString(w, "!type database\n")
  return nil
}

func (r *Repos) exportTableFromStruct(w io.Writer, tableName string, element interface{}) error {
  if err := r.exportTableHeaderFromStruct(w, tableName, element); err != nil {
    return err
  }
  return r.exportTableDataFromStruct(w, tableName, element)
}

func (r *Repos) exportTableHeaderFromStruct(w io.Writer, tableName string, element interface{}) error {
  io.WriteString(w, "\n!table " + tableName + "\n")
  colnames := `"` + strings.Join(structsql.ColumnNames(element), `","`) + `"`
  io.WriteString(w, "!columns " + colnames + "\n")
  return nil
}

func (r *Repos) exportTableDataFromStruct(w io.Writer, tableName string, element interface{}) error {
  sql, targets := structsql.SelectSql(tableName, element)
  sql = sql + ";"
  rows, err := r.db.Query(sql)
  if err != nil {
    return err
  }
  defer rows.Close()
  rowCount := 0
  for rows.Next() {
    err := rows.Scan(targets...)
    if err != nil {
      return err
    }
    // All of the targets are pointers to fields of the struct,
    // so if the struct field is a pointer (to allow null values),
    // then the type of the target is a double-pointer.
    for i, target := range targets {
      if i > 0 {
        io.WriteString(w, ",")
      }
      switch t := target.(type) {
      case *int:
        fmt.Fprintf(w, "%d", *t)
      case *string:
        fmt.Fprintf(w, "%q", *t)
      case **string:
        if *t == nil {
          fmt.Fprintf(w, "null")
        } else {
          fmt.Fprintf(w, "%q", **t)
        }
      default:
        // If we don't understand the type, print it out so that we know what
        // we need to add to the above switch statement.
        fmt.Fprintf(w, "(Type %T)", target)
      }
    }
    io.WriteString(w, "\n")
    rowCount++
  }
  if rowCount == 0 {
    io.WriteString(w, "#no rows\n")
  }
  return rows.Err()
}
