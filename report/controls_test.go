package report

import (
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestClientVisibleReports(t *testing.T) {
  reports, err := ClientVisibleReports([]string{"template"})
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
