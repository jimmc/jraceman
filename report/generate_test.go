package report

import (
  "io/ioutil"
  "testing"

  dbtest "github.com/jimmc/jracemango/dbrepo/test"
)

func TestStandardReports(t *testing.T) {
  for _, tt := range []struct{
    testName string
    reportName string
    data string
  } {
      { "test1", "test1", "<topdata>" },
      { "lanes", "lanes-test", "EV123" },
  } {
    t.Run(tt.testName, func(t *testing.T) {
      db, err := dbtest.DbWithTestTable()
      if err != nil {
        t.Fatal(err.Error())
      }
      defer db.Close()

      reportRoots := []string{"testdata", "form"}

      results, err := GenerateResults(db, reportRoots, tt.reportName, tt.data)
      if err != nil {
        t.Fatal(err)
      }
      outfile := "testdata/" + tt.reportName + ".out"
      err = ioutil.WriteFile(outfile, []byte(results.HTML), 0644)
      if err != nil {
        t.Fatal(err)
      }
      goldenfile := "testdata/" + tt.reportName + ".golden"
      if err := dbtest.CompareOutToGolden(outfile, goldenfile); err != nil {
        t.Fatal(err)
      }
    })
  }
}
