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
      Where: []ControlsWhereItem{
	{Name: "event_id", Display: "Event"},
	{Name: "event_name", Display: "Event Name"},
	{Name: "event_number", Display: "Event Number"},
	{Name: "team_id", Display: "Team"},
	{Name: "team_shortname", Display: "Team Short Name"},
	{Name: "team_name", Display: "Team Name"},
	{Name: "person_id", Display: "Person"},
      },
    },
    {
      Name: "org.jimmc.jraceman.Lanes",
      Display: "Lanes",
      Description: "",
      OrderBy: []ControlsOrderByItem{},
      Where: []ControlsWhereItem{},
    },
  }
  got, want := reports, expected
  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("ClientVisibleReports() mismatch (-want +got):\n%s", diff)
  }
}
