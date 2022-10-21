package report

import (
  "encoding/json"
  "fmt"
  "os"
  "testing"

  "github.com/jimmc/jraceman/dbrepo"

  goldenbase "github.com/jimmc/golden/base"
  _ "github.com/mattn/go-sqlite3"
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
  outfile := "testdata/client-visible-reports.out"
  goldenfile := "testdata/client-visible-reports.golden"
  f, err := os.Create(outfile)
  if err != nil {
    t.Fatalf("Error opening outfile: %v", err)
  }
  s, err := json.MarshalIndent(reports, "", "\t")
  if err != nil {
    t.Fatalf("Error encoding json: %v", err)
  }
  fmt.Fprintf(f, "%s", s)
  if err := f.Close(); err != nil {
    t.Fatalf("Error closing outfile: %v", err)
  }
  if err := goldenbase.CompareOutToGolden(outfile, goldenfile); err != nil {
    t.Fatalf("Error comparing to golden: %v", err)
  }
}
