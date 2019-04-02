package report

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestClientVisibleReports(t *testing.T) {
  reports, err := ClientVisibleReportsOne("template")
  if err != nil {
    t.Fatalf("ClientVisibleReports error: %v", err)
  }
  expected := []*ReportControls{
    {
      Name: "org.jimmc.jraceman.Entries",
      Display: "Entries",
      Description: "Entries ordered as selected.",
      OrderBy: []ControlsOrderByItem{
        ControlsOrderByItem{Name: "team",        Display: "Team, Person, Event"},
        ControlsOrderByItem{Name: "person",      Display: "Person, Event"},
        ControlsOrderByItem{Name: "eventTeam",   Display: "Event, Team, Person"},
        ControlsOrderByItem{Name: "eventPerson", Display: "Event, Person"},
      },
    },
    {
      Name: "org.jimmc.jraceman.Lanes",
      Display: "Lanes",
      Description: "",
      OrderBy: []ControlsOrderByItem{},
    },
  }
  got, want := reports, expected
  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("ClientVisibleReports() mismatch (-want +got):\n%s", diff)
  }
}

func TestExtractUserOrderByMap(t *testing.T) {
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
      got, err := extractUserOrderByList(tc.input)
      if tc.expectError {
        if err == nil {
          t.Fatalf("extractUserOrderByList: expected error but did not get one")
        }
      } else if err != nil {
        t.Fatalf("extractUserOrderByList: unexpected error: %v", err)
      } else {
        want := tc.expect
        if diff := cmp.Diff(want, got); diff != "" {
          t.Errorf("extractUserOrderByList mismatch (-want +got):\n%s", diff)
        }
      }
    })
  }
}
