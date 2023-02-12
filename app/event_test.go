package app

import (
  "encoding/json"
  "io/ioutil"
  "testing"

  dbtest "github.com/jimmc/jraceman/dbrepo/test"

  goldenbase "github.com/jimmc/golden/base"
)

func TestEventRaceInfo(t *testing.T) {
  tests := []struct{
    testName string
    setupName string
    outName string
    eventId string
    expectError bool
  } {
      { "no event id", "eventraceinfo", "", "", true },
      { "no such event", "eventraceinfo", "", "XYZ", true },
      { "no races", "eventraceinfo", "eventraceinfo-noraces", "M1.EV1", false },
      { "one race", "eventraceinfo", "eventraceinfo-onerace", "M1.EV3", false },
      { "more races", "eventraceinfo", "eventraceinfo-morefaces", "M1.EV4", false },
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {

      setupfilename := "testdata/" + tt.setupName + ".setup"

      dbRepos, err := dbtest.ReposAndLoadFile(setupfilename)
      if err != nil {
        t.Fatalf(err.Error())
      }
      defer dbRepos.Close()

      eventInfo, err := dbRepos.EventInfo().EventRaceInfo(tt.eventId)
      if tt.expectError {
        if err == nil {
          t.Fatal("Expected error but did not get one")
        }
      } else {
        if err != nil {
          t.Fatal(err)
        }

        jsonData, err := json.MarshalIndent(eventInfo, "", " ")
        if err != nil {
          t.Fatal(err)
        }
        outName := tt.outName
        if outName=="" {
          outName = tt.setupName
        }
        outfile := "testdata/" + outName + ".out"
        err = ioutil.WriteFile(outfile, jsonData, 0644)
        if err != nil {
          t.Fatal(err)
        }
        goldenfile := "testdata/" + outName + ".golden"
        if err := goldenbase.CompareOutToGolden(outfile, goldenfile); err != nil {
          t.Fatal(err)
        }
      }
    })
  }
}
