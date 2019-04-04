package report

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestExtractControlsOrderByList(t *testing.T) {
  tests := []struct{
    name string
    input *ReportAttributes
    expect []ControlsOrderByItem
    expectError bool
  } {
    {
      name: "empty",
      input: &ReportAttributes{},
      expect: []ControlsOrderByItem{},
    },
    {
      name: "normal",
      input: &ReportAttributes{
        Name: "normal",
        OrderBy: []AttributesOrderByItem{
          {Name:"a", Display: "AA"},
          {Name:"b", Display: "BB"},
        },
      },
      expect: []ControlsOrderByItem{
        {Name: "a", Display: "AA"},
        {Name: "b", Display: "BB"},
      },
    },
  }
  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      got, err := extractControlsOrderByItems(tc.input)
      if tc.expectError {
        if err == nil {
          t.Fatalf("extractControlsOrderByList: expected error but did not get one")
        }
      } else if err != nil {
        t.Fatalf("extractControlsOrderByList: unexpected error: %v", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got); diff != "" {
          t.Errorf("extractControlsOrderByList mismatch (-want +got):\n%s", diff)
        }
      }
    })
  }
}
