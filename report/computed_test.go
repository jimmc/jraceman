package report

import (
  "testing"
)

func TestGetComputed(t *testing.T) {
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
    { "nil attrs", &ReportOptions{OrderBy:"foo"}, nil, true },
    { "empty attrs map", &ReportOptions{OrderBy:"foo"}, attrsEmpty, true },
    { "invalid orderby option", &ReportOptions{OrderBy:"foo"}, attrsWithOrderBy, true },
    { "valid orderby option", &ReportOptions{OrderBy:"abc"}, attrsWithOrderBy, false },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      tplName := "<" + tc.name+ ">"
      computed, err := getComputed(tplName, tc.options, tc.attrs)
      if tc.expectError {
        if err == nil {
          t.Fatalf("Expected error but did not get it")
        }
      } else {
        if err != nil {
          t.Fatalf("Unexpected error: %v", err)
        }
        // TODO - we should check the computed result here.
        if computed == nil {
          t.Fatalf("Expected computed value but got nil")
        }
      }
    })
  }
}
