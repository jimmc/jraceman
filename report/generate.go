package report

import (
  "fmt"
  "strings"

  "github.com/jimmc/gtrepgen/dbsource"
  "github.com/jimmc/gtrepgen/gen"
)

type ReportOptions struct {
  OrderByKey string     // One of the key values in the attr orderby list for the template.
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

  attrs, err := gen.FindAndReadAttributes(templateName, reportRoots)
  if err != nil {
    return nil, fmt.Errorf("reading template attributes: %v", err)
  }

  attrsFunc := func(names ...string) (interface{}, error) {
    return descendAttributes(attrs, names...)
  }
  optionsFunc := func(names ...string) (interface{}, error) {
    return descendOptions(options, names...)
  }
  g := gen.New(templateName, true, w, dataSource)
  g = g.WithFuncs(map[string]interface{}{
    "attrs": attrsFunc,
    "options": optionsFunc,
  })
  if err := g.FromForm(reportRoots, data); err != nil {
    return nil, err
  }

  return &ReportResults{
    HTML: w.String(),
  }, nil
}

func descendAttributes(attrs interface{}, names ...string) (interface{}, error) {
  var a interface{}
  a = attrs
  for depth, name := range names {
    m, ok := a.(map[string]interface{})
    if !ok {
      return nil, fmt.Errorf("value is not map when trying to get field %q at depth %d", name, depth)
    }
    a, ok = m[name]
    if !ok {
      return nil, fmt.Errorf("field %q not found at depth %d", name, depth)
    }
  }
  return a, nil
}

func descendOptions(options *ReportOptions, names ...string) (interface{}, error) {
  if options == nil {
    return nil, fmt.Errorf("no report options specified")
  }
  if len(names) == 0 {
    return nil, fmt.Errorf("no option name specified")
  }
  if len(names) > 1 {
    return nil, fmt.Errorf("too many option names specified")
  }
  switch names[0] {
  case "orderby":
    return options.OrderByKey, nil
  default:
    return nil, fmt.Errorf("unknown option name %q", names[0])
  }
}
