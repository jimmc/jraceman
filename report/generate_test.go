package report

import (
  "io/ioutil"
  "testing"

  dbtest "github.com/jimmc/jracemango/dbrepo/test"
)

func TestReportBasic(t *testing.T) {
  db, err := dbtest.DbWithTestTable()
  if err != nil {
    t.Fatal(err.Error())
  }
  defer db.Close()

  reportRoot := "testdata"
  name := "test1"
  data := "<topdata>"

  results, err := GenerateResults(db, reportRoot, name, data)
  if err != nil {
    t.Fatal(err)
  }
  outfile := "testdata/" + name + ".out"
  err = ioutil.WriteFile(outfile, []byte(results.HTML), 0644)
  if err != nil {
    t.Fatal(err)
  }
  goldenfile := "testdata/" + name + ".golden"
  if err := dbtest.CompareOutToGolden(outfile, goldenfile); err != nil {
    t.Fatal(err)
  }
}
