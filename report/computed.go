package report

import (
  "fmt"
)

// ReportComputed is the collection of data that we compute based on the report attributes
// and the ReportOptions.
type ReportComputed struct {
  OrderByExpr string    // The orderby sql corresponding to the OrderByKey in the Options.
  OrderByDisplay string
  OrderByClause string  // Blank if no OrderBy, else "ORDER BY " and the expression.
}

func validateReportOptions(tplName string, options *ReportOptions, attrs *ReportAttributes, computed *ReportComputed) error {
  if options == nil {
    return nil
  }
  if options.OrderByKey != "" {
    if attrs == nil || len(attrs.OrderBy) == 0 {
      return fmt.Errorf("invalid orderby option %q, template %s does not permit orderby",
          options.OrderByKey, tplName)
    }
    orderByItem, err := findOrderByItem(attrs.OrderBy, options.OrderByKey)
    if err != nil {
      return fmt.Errorf("invalid orderby option %q for template %s",
          options.OrderByKey, tplName)
    }
    computed.OrderByExpr = orderByItem.Sql
    computed.OrderByDisplay = orderByItem.Display
    if computed.OrderByExpr != "" {
      computed.OrderByClause = "ORDER BY " + computed.OrderByExpr
    }
  }
  return nil
}

func findOrderByItem(orderByList []AttributesOrderByItem, orderByName string) (*AttributesOrderByItem, error) {
  for _, item := range orderByList {
    if item.Name == orderByName {
      return &item, nil
    }
  }
  return nil, fmt.Errorf("orderby item %q not found", orderByName)
}
