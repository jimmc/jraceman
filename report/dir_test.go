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
    {Name: "org.jimmc.jraceman.Entries", Display: "Entries" },
    {Name: "org.jimmc.jraceman.Lanes", Display: "Lanes" },
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
