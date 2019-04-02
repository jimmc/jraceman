package report

import (
  "fmt"
  "testing"
)

func TestReadTemplateAttrs(t *testing.T) {
  attrslist, err := ReadTemplateAttrs("template")
  if err != nil {
    t.Fatalf("ReadTemplateAttrs error: %v", err)
  }
  if got, want := len(attrslist), 2; got != want {
    t.Fatalf("Wrong number of files-with-attributes found, got %d, want %d", got, want)
  }
  fmt.Printf("attrslist: %+v\n", attrslist)
  if got, want := attrslist[0].Name, "org.jimmc.jraceman.Entries"; got != want {
    t.Errorf("Name of first report: got %s, want %s", got, want)
  }
}
