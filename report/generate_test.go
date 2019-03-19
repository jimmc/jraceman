package report

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "testing"
  "time"

  "github.com/jimmc/gtrepgen/gen"
  "github.com/jimmc/jracemango/dbrepo"
  dbtest "github.com/jimmc/jracemango/dbrepo/test"
)

func tNow() time.Time {
  format := "Mon Jan 2 2006 15:04:05"
  testTime := "Tue Mar 12 2018 10:11:12"
  t, err := time.Parse(format, testTime)
  if err != nil {
    panic("Failed to parse test time: " + err.Error())
  }
  return t
}

func TestStandardReports(t *testing.T) {
  oldNowFunc := gen.Now
  gen.Now = tNow
  defer func() { gen.Now = oldNowFunc }()

  roots1 := []string{"template"}
  roots2 := []string{"testdata", "template"}
  for _, tt := range []struct{
    testName string
    setupName string
    reportName string
    reportRoots []string
    templateName string
    data string
  } {
      { "test1", "empty", "test1", roots2, "test1", "<topdata>" },
      { "lanes", "empty", "lanes-test", roots2, "lanes-test", "EV123" },
      { "entries", "sample1", "entries-test", roots1, "org.jimmc.jraceman.Entries", "M1.EV1" },
  } {
    t.Run(tt.testName, func(t *testing.T) {

      setupfilename := "testdata/" + tt.setupName + ".setup"

      dbRepos, err := openAndLoad(setupfilename)
      if err != nil {
        t.Fatalf(err.Error())
      }
      defer dbRepos.Close()
      db := dbRepos.DB()

      reportOpts := &ReportOptions{}
      results, err := GenerateResults(db, tt.reportRoots, tt.templateName, tt.data, reportOpts)
      if err != nil {
        t.Fatal(err)
      }
      outfile := "testdata/" + tt.reportName + ".html"
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

func openAndLoad(setupfile string) (*dbrepo.Repos, error) {
  dbRepos, err := dbrepo.Open("sqlite3::memory:")
  if err != nil {
    return nil, fmt.Errorf("failed to open repository: %v", err)
  }

  err = dbRepos.CreateTables()
  if err != nil {
    return nil, fmt.Errorf("failed to create repository tables: %v", err)
  }

  inFile, err := os.Open(setupfile)
  if err != nil {
    return nil, fmt.Errorf("error opening import input file %s: %v", setupfile, err)
  }
  defer inFile.Close()

  log.Printf("Importing from %s\n", setupfile)
  counts, err := dbRepos.Import(inFile)
  if err != nil {
    return nil, fmt.Errorf("error importing from %s: %v", setupfile, err)
  }
  log.Printf("Import done: inserted %d, updated %d, unchanged %d records\n",
      counts.Inserted(), counts.Updated(), counts.Unchanged())
  return dbRepos, nil
}
