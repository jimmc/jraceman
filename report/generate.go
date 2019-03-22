package report

import (
  "fmt"
  "strings"

  "github.com/jimmc/gtrepgen/dbsource"
  "github.com/jimmc/gtrepgen/gen"
)

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

  attrs, err := gen.FindAndReadAttributes(templateName, reportRoots)
  if err != nil {
    return nil, fmt.Errorf("reading template attributes: %v", err)
  }
  attrsMap := map[string]interface{}{}
  if attrs != nil {
    var ok bool
    attrsMap, ok = attrs.(map[string]interface{})
    if !ok {
      return nil, fmt.Errorf("attributes in template %q is not a map", templateName)
    }
  }

  if options != nil {
    if err := validateReportOptions(templateName, options, attrsMap); err != nil {
      return nil, err
    }
  }

  attrsFunc := func(names ...string) (interface{}, error) {
    return descendAttributes(attrsMap, names...)
  }
  optionsFunc := func(names ...string) (interface{}, error) {
    return descendOptions(options, names...)
  }
  whereData, err := where(attrsMap, options)
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
    "where": whereFunc,
  })
  if err := g.FromTemplate(reportRoots, data); err != nil {
    return nil, err
  }

  return &ReportResults{
    HTML: w.String(),
  }, nil
}

func validateReportOptions(tplName string, options *ReportOptions, attrsMap map[string]interface{}) error {
  if options == nil {
    return nil
  }
  if options.OrderByKey != "" {
    attrsOrderby := attrsMap["orderby"]
    if attrsOrderby == nil {
      return fmt.Errorf("invalid orderby option %q, template %s does not permit orderby",
          options.OrderByKey, tplName)
    }
    orderbyMap, ok := attrsOrderby.(map[string]interface{})
    if !ok {
      return fmt.Errorf("invalid format for orderby map in template %s", tplName)
    }
    if orderbyMap[options.OrderByKey] == nil {
      return fmt.Errorf("invalid orderby option %q for template %s",
          options.OrderByKey, tplName)
    }
  }
  return nil
}

func descendAttributes(attrsMap map[string]interface{}, names ...string) (interface{}, error) {
  var a interface{}
  a = attrsMap
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
