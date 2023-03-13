package dbrepo_test

import (
  "context"
  "encoding/json"
  "io/ioutil"
  "os"
  "testing"

  "github.com/jimmc/jraceman/domain"
  "github.com/jimmc/jraceman/dbrepo"
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
      { "more races", "eventraceinfo", "eventraceinfo-moreraces", "M1.EV4", false },
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {

      // Load the database.
      setupfilename := "testdata/" + tt.setupName + ".setup"
      dbRepos, err := dbtest.ReposAndLoadFile(setupfilename)
      if err != nil {
        t.Fatalf(err.Error())
      }
      defer dbRepos.Close()

      // Run the function under test.
      eventInfo, err := dbRepos.EventInfo().EventRaceInfo(tt.eventId)

      // Check the result.
      if tt.expectError {
        if err == nil {
          t.Fatal("Expected error but did not get one")
        }
      } else {
        if err != nil {
          t.Fatal(err)
        }

        // Write out the database.
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

        // Check that the result is as expected.
        goldenfile := "testdata/" + outName + ".golden"
        if err := goldenbase.CompareOutToGolden(outfile, goldenfile); err != nil {
          t.Fatal(err)
        }
      }
    })
  }
}

// raceInfoA is a race that does not exist in updateraceinfo.setup.
var raceInfoA = &domain.RaceInfo{
  RaceID: "R7",
  EventID: "M1.EV4",
  StageID: "S1",
  AreaID: "",
  StageName: "Heat",
  StageNumber: 1,
  IsFinal: false,
  Round: 1,
  Section: 3,
  AreaName: "",
  RaceNumber: 1007,
  LaneCount: 5,
}
// raceInfoB is a race that does exist in updateraceinfo.setup.
var raceInfoB = &domain.RaceInfo{
  RaceID: "R5",
  EventID: "M1.EV4",
  StageID: "S2",
  AreaID: "",
  StageName: "Semi",
  StageNumber: 2,
  IsFinal: false,
  Round: 2,
  Section: 2,
  AreaName: "",
  RaceNumber: 1006,
  LaneCount: 0,
}
// raceInfoB2 is a race that does exist in updateraceinfo.setup, but modified.
var raceInfoB2 = &domain.RaceInfo{
  RaceID: "R5",
  EventID: "M1.EV4",
  StageID: "S1",
  AreaID: "A1",
  StageName: "Heat",
  StageNumber: 1,
  IsFinal: false,
  Round: 1,
  Section: 4,
  AreaName: "",
  RaceNumber: 1006,
  LaneCount: 0,
}
// raceInfoBad is a race that does contain sufficient information for creation.
var raceInfoBad = &domain.RaceInfo{
  RaceID: "R1",
  EventID: "",
  StageID: "",
  AreaID: "",
  StageName: "",
  StageNumber: 0,
  IsFinal: false,
  Round: 0,
  Section: 0,
  AreaName: "",
  RaceNumber: 0,
  LaneCount: 0,
}
var raceInfoListA = []*domain.RaceInfo{ raceInfoA }
var raceInfoListB = []*domain.RaceInfo{ raceInfoB }
var raceInfoListB2 = []*domain.RaceInfo{ raceInfoB2 }
var raceInfoListBad = []*domain.RaceInfo{ raceInfoBad }

func TestUpdateRaceInfo(t *testing.T) {
  tests := []struct{
    testName string
    setupName string
    outName string
    expectError bool
    eventInfo *domain.EventInfo
    racesToCreate []*domain.RaceInfo
    racesToDelete []*domain.RaceInfo
    racesToModFrom []*domain.RaceInfo
    racesToModTo []*domain.RaceInfo
  } {
      { "no changes", "updateraceinfo", "updateraceinfo-nochange", false, nil, nil, nil, nil, nil },
      { "add race", "updateraceinfo", "updateraceinfo-addrace", false, nil, raceInfoListA, nil, nil, nil },
      { "delete race", "updateraceinfo", "updateraceinfo-deleterace", false, nil, nil, raceInfoListB, nil, nil },
      { "modify race", "updateraceinfo", "updateraceinfo-modifyrace", false, nil, nil, nil, raceInfoListB, raceInfoListB2 },
      { "add bad race", "updateraceinfo", "updateraceinfo-nochange", true, nil, raceInfoListBad, nil, nil, nil },
      { "delete bad race", "updateraceinfo", "updateraceinfo-nochange", true, nil, nil, raceInfoListA, nil, nil },
      { "modify bad race", "updateraceinfo", "updateraceinfo-nochange", true, nil, nil, nil, raceInfoListBad, raceInfoListB2 },
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {

      // Load the database.
      setupfilename := "testdata/" + tt.setupName + ".setup"
      dbRepos, err := dbtest.ReposAndLoadFile(setupfilename)
      if err != nil {
        t.Fatalf(err.Error())
      }
      defer dbRepos.Close()

      // Run the function under test.
      ctx := context.Background()
      err = dbRepos.EventInfo().UpdateRaceInfo(ctx, tt.eventInfo,
        tt.racesToCreate, tt.racesToDelete, tt.racesToModFrom, tt.racesToModTo)

      // Check the result.
      if tt.expectError {
        if err == nil {
          t.Fatal("Expected error but did not get one")
        }
      } else {
        if err != nil {
          t.Fatal(err)
        }
      }

      // We write out our data and compare to our golden even if
      // we were expecting an error. If we got the expected error,
      // we want to check that the database did not change.

      // Write out the Race table.
      outName := tt.outName
      if outName=="" {
        outName = tt.setupName
      }
      outfile := "testdata/" + outName + ".out"
      w, err := os.Create(outfile)
      if err != nil {
        t.Fatalf("error opening export output file %s: %v", outfile, err)
      }
      exporter, err := dbRepos.NewExporter()
      if err != nil {
        t.Fatal(err)
      }
      dbReposRace := dbRepos.Race().(*dbrepo.DBRaceRepo)
      if err := dbReposRace.Export(exporter, w); err != nil {
        t.Fatalf("error exporting to %s: %v", outfile, err)
      }
      w.Close()

      // Check that the result is as expected.
      goldenfile := "testdata/" + outName + ".golden"
      if err := goldenbase.CompareOutToGolden(outfile, goldenfile); err != nil {
        t.Fatal(err)
      }
    })
  }
}
