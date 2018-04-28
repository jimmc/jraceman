package dbrepo

import (
  "io"
  "strings"
)

func (r *Repos) Export(w io.Writer) error {
  if err := r.exportHeader(w); err != nil {
    return err
  }
  //r.Area().Export(w)
  r.dbSite.Export(r, w)
  return nil
}

func (r *Repos) exportHeader(w io.Writer) error {
  io.WriteString(w, "#!jraceman -import\n")
  io.WriteString(w, "!exportVersion 2\n")
  io.WriteString(w, "!appInfo JRaceman v2.0.0\n")       // TODO - get real version
  io.WriteString(w, "!type database\n")
  return nil
}

func (r *Repos) exportTableHeaderFromStruct(w io.Writer, tableName string, element interface{}) error {
  io.WriteString(w, "\n!table " + tableName + "\n")
  colnames := `"` + strings.Join(stdColumnNamesFromStruct(element), `","`) + `"`
  io.WriteString(w, "!columns " + colnames + "\n")
  return nil
}
