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

func getComputed(templateName string, options *ReportOptions, attrs *ReportAttributes) (*ReportComputed, error) {
  if attrs == nil {
    return nil, fmt.Errorf("nil attrs are not allowed in getComputed") // Internal error.
  }
  where, err := computeWhere(attrs, options)
  if err != nil {
    return nil, err
  }
  orderby, err := computeOrderBy(templateName, options, attrs)
  if err != nil {
    return nil, err
  }
  computed := &ReportComputed{
    Where: *where,
    OrderBy: *orderby,
  }
  return computed, nil
}
