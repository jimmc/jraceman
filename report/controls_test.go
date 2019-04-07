package report

import (
  "testing"

  "github.com/jimmc/jracemango/dbrepo"

  "github.com/google/go-cmp/cmp"
)

func TestClientVisibleReports(t *testing.T) {
  dbrepos, err := dbrepo.Open("sqlite3::memory:")
  if err != nil {
    t.Fatalf("failed to open repository: %v", err)
  }
  defer dbrepos.Close()
  reports, err := ClientVisibleReports(dbrepos, []string{"template"})
  if err != nil {
    t.Fatalf("ClientVisibleReports error: %v", err)
  }
  keyOps := []string{"eq"}
  stringOps := []string{"eq", "ne", "gt", "ge", "lt", "le", "like"}
  dfltOps := []string{"eq", "ne", "gt", "ge", "lt", "le"}
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
	{Name: "event_id", Display: "Event", Ops: keyOps},
	{Name: "event_name", Display: "Event Name", Ops: stringOps},
	{Name: "event_number", Display: "Event Number", Ops: dfltOps},
	{Name: "team_id", Display: "Team", Ops: keyOps},
	{Name: "team_shortname", Display: "Team Short Name", Ops: stringOps},
	{Name: "team_name", Display: "Team Name", Ops: stringOps},
	{Name: "person_id", Display: "Person", Ops: keyOps},
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
