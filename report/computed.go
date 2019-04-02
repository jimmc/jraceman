package report

import (
  "fmt"
)

// ReportComputed is the collection of data that we compute based on the report attributes
// and the ReportOptions.
type ReportComputed struct {
  OrderBy ComputedOrderBy
  Where ComputedWhere
}

type ComputedOrderBy struct {
  Expr string    // The orderby sql corresponding to the OrderByKey in the Options.
  Display string
  Clause string  // Blank if no OrderBy, else "ORDER BY " and the expression.
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
    computed.OrderBy.Expr = orderByItem.Sql
    computed.OrderBy.Display = orderByItem.Display
    if computed.OrderBy.Expr != "" {
      computed.OrderBy.Clause = "ORDER BY " + computed.OrderBy.Expr
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

func getComputed(templateName string, options *ReportOptions, attrs *ReportAttributes) (*ReportComputed, error) {
  where, err := where(attrs, options)
  if err != nil {
    return nil, err
  }
  computed := &ReportComputed{
    Where: *where,
  }
  if options != nil {
    if err := validateReportOptions(templateName, options, attrs, computed); err != nil {
      return nil, err
    }
  }
  return computed, nil
}
