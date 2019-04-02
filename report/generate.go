package report

import (
  "strings"

  "github.com/jimmc/gtrepgen/dbsource"
  "github.com/jimmc/gtrepgen/gen"

  "github.com/golang/glog"
)

// ReportOptions is the data given to us by the user to generate an instance of the report.
type ReportOptions struct {
  OrderByKey string     // One of the key values in the attr orderby list for the template.
  WhereValues map[string]OptionsWhereValue
}

type ReportResults struct {
  HTML string           // The generated report.
}

// GenerateResults generates a report from a template and a top piece of data.
// The db argument can be either a sql.DB or a sql.Tx.
func GenerateResults(db dbsource.DBQuery, reportRoots []string, templateName, data string,
    options *ReportOptions) (*ReportResults, error) {
  dataSource := dbsource.New(db)
  w := &strings.Builder{}

  attrs, err := getAttributes(templateName, reportRoots)
  if err != nil {
    return nil, err
  }
  glog.V(1).Infof("attrs=%+v\n", attrs)
  computed, err := getComputed(templateName, options, attrs)
  if err != nil {
    return nil, err
  }
  glog.V(1).Infof("computed=%+v\n", computed)
  whereData, err := where(attrs, options)
  if err != nil {
    return nil, err
  }
  glog.V(1).Infof("where=%+v\n", whereData)

  g := gen.New(templateName, true, w, dataSource)
  g = g.WithFuncs(map[string]interface{}{
    "attrs": func() interface{} { return attrs },
    "options": func() interface{} { return options },
    "computed": func() interface{} { return computed },
    "where": func() interface{} { return whereData },
  })
  if err := g.FromTemplate(reportRoots, data); err != nil {
    return nil, err
  }

  return &ReportResults{
    HTML: w.String(),
  }, nil
}
