package report

import (
  "strings"

  "github.com/jimmc/gtrepgen/dbsource"
  "github.com/jimmc/gtrepgen/gen"
)

type ReportResults struct {
  HTML string           // The generated report.
}

// GenerateResults generates a report from a form and a top piece of data.
// The db argument can be either a sql.DB or a sql.Tx.
func GenerateResults(db dbsource.DBQuery, reportRoot, formname, data string) (*ReportResults, error) {
  dataSource := dbsource.New(db)
  refdirpath := reportRoot
  w := &strings.Builder{}

  g := gen.New(formname, true, w, dataSource)
  if err := g.FromForm(refdirpath, data); err != nil {
    return nil, err
  }

  return &ReportResults{
    HTML: w.String(),
  }, nil
}
