package report

import (
  "fmt"
  "strings"

  "github.com/jimmc/gtrepgen/dbsource"
  "github.com/jimmc/gtrepgen/gen"
)

type ReportResults struct {
  HTML string           // The generated report.
}

// GenerateResults generates a report from a template and a top piece of data.
// The db argument can be either a sql.DB or a sql.Tx.
func GenerateResults(db dbsource.DBQuery, reportRoots []string, templateName, data string) (*ReportResults, error) {
  dataSource := dbsource.New(db)
  w := &strings.Builder{}

  attrs, err := gen.FindAndReadAttributes(templateName, reportRoots)
  if err != nil {
    return nil, fmt.Errorf("reading template attributes: %v", err)
  }

  attrsFunc := func(names ...string) (interface{}, error) {
    return descendAttributes(attrs, names...)
  }
  g := gen.New(templateName, true, w, dataSource)
  g = g.WithFuncs(map[string]interface{}{
    "attrs": attrsFunc,
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
  for _, name := range names {
    m, ok := a.(map[string]interface{})
    if !ok {
      return nil, fmt.Errorf("value is not map when trying to get field %q", name)
    }
    a, ok = m[name]
    if !ok {
      return nil, fmt.Errorf("field %q not found", name)
    }
  }
  return a, nil
}
