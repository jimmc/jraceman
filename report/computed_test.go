package report

import (
  "testing"
)

func TestValidateReportOptions(t *testing.T) {
  attrsEmpty := &ReportAttributes{}
  attrsWithOrderBy := &ReportAttributes{
    OrderBy: []AttributesOrderByItem{
      AttributesOrderByItem{Name: "abc", Display: "xyz"},
      AttributesOrderByItem{Name: "def", Display: "uvw"},
    },
  }
  tests := []struct{
    name string
    options *ReportOptions
    attrs *ReportAttributes
    expectError bool
  }{
    { "nil options", nil, attrsEmpty, false },
    { "no orderby option", &ReportOptions{}, attrsEmpty, false },
    { "nil attrs", &ReportOptions{OrderByKey:"foo"}, nil, true },
    { "empty attrs map", &ReportOptions{OrderByKey:"foo"}, attrsEmpty, true },
    { "invalid orderby option", &ReportOptions{OrderByKey:"foo"}, attrsWithOrderBy, true },
    { "valid orderby option", &ReportOptions{OrderByKey:"abc"}, attrsWithOrderBy, false },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      tplName := "<" + tc.name+ ">"
      computed := &ReportComputed{}
      err := validateReportOptions(tplName, tc.options, tc.attrs, computed)
      if tc.expectError && err == nil {
        t.Errorf("Expected error but dit not get it")
      } else if !tc.expectError && err != nil {
        t.Errorf("Unexpected error: %v", err)
      }
    })
  }
}
