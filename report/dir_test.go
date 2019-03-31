package report

import (
  "fmt"
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestClientVisibleReports(t *testing.T) {
  reports, err := ClientVisibleReportsOne("template")
  if err != nil {
    t.Fatalf("ClientVisibleReports error: %v", err)
  }
  expected := []*ReportAttributes{
    {
      Name: "org.jimmc.jraceman.Entries",
      Display: "Entries",
      OrderBy: []OrderByItem{
        OrderByItem{Name: "team",        Display: "Team, Person, Event"},
        OrderByItem{Name: "person",      Display: "Person, Event"},
        OrderByItem{Name: "eventTeam",   Display: "Event, Team, Person"},
        OrderByItem{Name: "eventPerson", Display: "Event, Person"},
      },
    },
    {
      Name: "org.jimmc.jraceman.Lanes",
      Display: "Lanes",
    },
  }
  got, want := reports, expected
  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("ClientVisibleReports() mismatch (-want +got):\n%s", diff)
  }
}

func TestReadTemplateAttrs(t *testing.T) {
  attrslist, err := ReadTemplateAttrs("template")
  if err != nil {
    t.Fatalf("ReadTemplateAttrs error: %v", err)
  }
  if got, want := len(attrslist), 2; got != want {
    t.Fatalf("Wrong number of files-with-attributes found, got %d, want %d", got, want)
  }
  fmt.Printf("attrslist: %+v\n", attrslist)
  if got, want := attrslist[0]["name"], "org.jimmc.jraceman.Entries"; got != want {
    t.Errorf("Name of first report: got %s, want %s", got, want)
  }
}

func TestExtractUserOrderByMap(t *testing.T) {
  tests := []struct{
    name string
    input map[string]interface{}
    expect []OrderByItem
    expectError bool
  } {
    {
      name: "empty",
      input: map[string]interface{}{},
      expect: nil,
    },
    {
      name: "normal",
      input: map[string]interface{}{
        "name": "normal",
        "orderby": []interface{}{
          map[string]interface{}{"name":"a", "display": "AA"},
          map[string]interface{}{"name":"b", "display": "BB"},
        },
      },
      expect: []OrderByItem{
        {Name: "a", Display: "AA"},
        {Name: "b", Display: "BB"},
      },
    },
    {
      name: "bad_orderby",
      input: map[string]interface{}{
        "name": "bad_orderby",
        "orderby": "not_an_array",
      },
      expectError: true,
    },
    {
      name: "bad_orderby_item",
      input: map[string]interface{}{
        "name": "bad_orderby_item",
        "orderby": []interface{}{
          "not_a_map",
        },
      },
      expectError: true,
    },
    {
      name: "bad_orderby_display",
      input: map[string]interface{}{
        "name": "bad_orderby_display",
        "orderby": []interface{}{
          map[string]interface{}{"display": 123},
        },
      },
      expectError: true,
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
