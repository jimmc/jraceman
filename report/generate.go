package report

import (
  "fmt"
  "strings"

  "github.com/jimmc/gtrepgen/dbsource"
  "github.com/jimmc/gtrepgen/gen"

  "github.com/golang/glog"
)

// ReportOptions is the data given to us by the user to generate an instance of the report.
type ReportOptions struct {
  OrderByKey string     // One of the key values in the attr orderby list for the template.
  WhereValues map[string]WhereValue
}

// WhereValues contains the values specified in the options for one where field.
type WhereValue struct {
  Op string     // The comparison operation to use for this field.
  Value interface{}     // The value to use on the RHS of the comparison.
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

  attrs := &ReportAttributes{}
  if err := gen.FindAndReadAttributesInto(templateName, reportRoots, attrs); err != nil {
    return nil, fmt.Errorf("reading template attributes: %v", err)
  }
  computed := &ReportComputed{}
  if options != nil {
    if err := validateReportOptions(templateName, options, attrs, computed); err != nil {
      return nil, err
    }
  }
  glog.V(1).Infof("computed=%+v\n", computed)

  attrsFunc := func(names ...string) (interface{}, error) {
    return attrs, nil
  }
  optionsFunc := func(names ...string) (interface{}, error) {
    return options, nil
  }
  computedFunc := func(names ...string) (interface{}, error) {
    return computed, nil
  }
  whereData, err := where(attrs, options)
  if err != nil {
    return nil, err
  }
  whereFunc := func() (interface{}, error) {
    return whereData, nil
  }
  g := gen.New(templateName, true, w, dataSource)
  g = g.WithFuncs(map[string]interface{}{
    "attrs": attrsFunc,
    "options": optionsFunc,
    "computed": computedFunc,
    "where": whereFunc,
  })
  if err := g.FromTemplate(reportRoots, data); err != nil {
    return nil, err
  }

  return &ReportResults{
    HTML: w.String(),
  }, nil
}
