package report

import (
)

// ReportComputed is the collection of data that we compute based on the report attributes
// and the ReportOptions.
type ReportComputed struct {
  OrderBy ComputedOrderBy
  Where ComputedWhere
}

func validateReportOptions(tplName string, options *ReportOptions, attrs *ReportAttributes, computed *ReportComputed) error {
  if options == nil {
    return nil
  }
  if options.OrderByKey != "" {
    ob, err := computeOrderBy(tplName, options, attrs)
    if err != nil {
      return err
    }
    computed.OrderBy = *ob
  }
  return nil
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
