package app

import (
  "context"
  "encoding/json"
  "io/ioutil"
  "errors"
  "reflect"
  "testing"

  "github.com/jimmc/jraceman/domain"
  dbtest "github.com/jimmc/jraceman/dbrepo/test"

  goldenbase "github.com/jimmc/golden/base"
)

func TestEventCreateRaces(t *testing.T) {
  tests := []struct{
    testName string
    setupName string
    outName string
    eventId string
    laneCount int
    dryRun bool
    allowDeleteLanes bool
    expectError bool
  } {
      { "no event id", "eventcreateraces-errors", "", "", 0, false, false, true },
      { "no such event", "eventcreateraces-errors", "", "XYZ", 0, false, false, true },
      { "no progression", "eventcreateraces-errors", "", "M1.EV2", 0, false, false, true },
      { "no races no entries", "eventcreateraces-errors", "", "M1.EV4", 0, false, false, true },
      //{ "create one race", "eventcreateraces", "eventcreateraces-onerace", "M1.EV3", 0, false, false, false },
      //{ "more races", "eventcreateraces", "eventcreateraces-moreraces", "M1.EV4", 0, false, false, false },
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

      ctx := context.Background()
      // racesResult type is *CreateRacesResult
      // Run the function under test.
      racesResult, err := EventCreateRaces(ctx, dbRepos, tt.eventId, tt.laneCount,
            tt.dryRun, tt.allowDeleteLanes)

      // Check the result.
      if tt.expectError {
        if err == nil {
          t.Fatal("Expected error but did not get one")
        }
      } else {
        if err != nil {
          t.Fatal(err)
        }

        // Write out the result.
        jsonData, err := json.MarshalIndent(racesResult, "", " ")
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

var raceInfoOneHeat = &domain.RaceInfo{
  Round: 1,
  Section: 1,
  StageNumber: 1,
  LaneCount: 0,
}
var raceInfoOneDirectFinal = &domain.RaceInfo{
  IsFinal: true,
  Round: 1,
  Section: 1,
  StageNumber: 3,
  LaneCount: 0,
}
var raceInfoOneDirectFinalWithLanes = &domain.RaceInfo{
  IsFinal: true,
  Round: 1,
  Section: 1,
  StageNumber: 3,
  LaneCount: 5,
}
var raceInfosEmpty = []*domain.RaceInfo{}
var raceInfosOneHeat = []*domain.RaceInfo{ raceInfoOneHeat }
var raceInfosOneDirectFinal = []*domain.RaceInfo{ raceInfoOneDirectFinal }
var raceInfosOneDirectFinalWithLanes = []*domain.RaceInfo{ raceInfoOneDirectFinalWithLanes }

var eventInfoEmpty = &domain.EventInfo{ }

var eventInfoNoRaces = &domain.EventInfo{
  EventID: "E1",
  Summary: "E1: No Entries",
  RoundCounts: make([]*domain.EventRoundCounts, 0),
  Races: raceInfosEmpty,
}

var eventInfoNoRacesWithEntryCount = &domain.EventInfo{
  EventID: "E1",
  Summary: "E1: No Entries",
  RoundCounts: make([]*domain.EventRoundCounts, 0),
  EntryCount: 5,
  Races: raceInfosEmpty,
}

var eventInfoNoRacesWithGroupCount = &domain.EventInfo{
  EventID: "E1",
  Summary: "E1: No Entries",
  RoundCounts: make([]*domain.EventRoundCounts, 0),
  EntryCount: 13,
  GroupSize: 2,
  GroupCount: 6,
  Races: raceInfosEmpty,
}

var eventInfoOneDirectFinal = &domain.EventInfo{
  EventID: "E1",
  Summary: "E1: No Entries",
  RoundCounts: make([]*domain.EventRoundCounts, 0),
  Races: raceInfosOneDirectFinal,
}

var eventRoundCountsEmpty = []*domain.EventRoundCounts{}

var eventRoundCountsOneFinal = []*domain.EventRoundCounts{
  {
    Count: 1,
    Round: 1,
    StageNumber: 3,
    IsFinal: true,
  },
}

// progSysT allows us to control what gets returned by the ProgSys
// we use for testing.
type progSysT struct{
  roundCounts []*domain.EventRoundCounts
  err error
}
func (p *progSysT) DesiredRoundCounts(progressionState string, laneCount, areaLanes, areaExtraLanes int) ([]*domain.EventRoundCounts, error) {
  return p.roundCounts, p.err
}

var progSysErr = &progSysT{nil, errors.New("Intentional error for testing")}
var progSysNoRounds = &progSysT{make([]*domain.EventRoundCounts,0), nil}
var progSysOneFinal = &progSysT{eventRoundCountsOneFinal, nil}

func TestCalculateRaceChanges(t *testing.T) {
  tests := []struct{
    testName string
    event *domain.EventInfo
    progression ProgSys
    laneCount int
    expectError bool
    createCount int
    deleteCount int
    modifyCount int
  } {
    { "empty event" , eventInfoEmpty, nil, 0, true, 0, 0, 0 },
    { "no races and no lanes" , eventInfoNoRaces, nil, 0, true, 0, 0, 0 },
    { "one final no changes" , eventInfoOneDirectFinal, progSysOneFinal, 5, false, 0, 0, 0 },
    { "add one final", eventInfoNoRacesWithEntryCount, progSysOneFinal, -1, false, 1, 0, 0 },
    { "add one group final", eventInfoNoRacesWithGroupCount, progSysOneFinal, -1, false, 1, 0, 0 },
    { "progression error", eventInfoNoRacesWithGroupCount, progSysErr, -1, true, 0, 0, 0 },
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {

      result, err := CalculateRaceChangesForTesting(tt.event, tt.progression, tt.laneCount)

      if tt.expectError {
        if err == nil {
          t.Fatal("Expected error but did not get one")
        }
        return
      }
      if err != nil {
        t.Fatal(err)
      }
      if result == nil {
        t.Fatal("Expected result but did not get it")
      }
      if got, want := len(result.RacesToCreate), tt.createCount; got!=want {
        t.Errorf("Count of races to create: got %d, want %d", got, want)
      }
      if got, want := len(result.RacesToDelete), tt.deleteCount; got!=want {
        t.Errorf("Count of races to delete: got %d, want %d", got, want)
      }
      if got, want := len(result.RacesToModFrom), tt.modifyCount; got!=want {
        t.Errorf("Count of races to mod from: got %d, want %d; got modfrom=%+v, modto=%+v",
            got, want, result.RacesToModFrom, result.RacesToModTo)
      }
      if len(result.RacesToModFrom)!=len(result.RacesToModTo) {
        t.Errorf("Expected len(modFrom)==len(modTo), but got %d and %d",
            len(result.RacesToModFrom), len(result.RacesToModTo))
      }
    })
  }
}

func TestFirstRaceWithLaneData(t *testing.T) {
  tests := []struct{
    testName string
    races []*domain.RaceInfo
    result *domain.RaceInfo
  } {
    { "no races", raceInfosEmpty, nil },
    { "one race without lanes", raceInfosOneDirectFinal, nil },
    { "one race with lanes", raceInfosOneDirectFinalWithLanes, raceInfoOneDirectFinalWithLanes },
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {
      if got, want := firstRaceWithLaneData(tt.races), tt.result; !reflect.DeepEqual(got, want) {
        t.Errorf("raceInfo got %v, want %v", got, want)
      }
    })
  }
}

func TestRoundsToRaces(t *testing.T) {
  tests := []struct{
    testName string
    roundCounts []*domain.EventRoundCounts
    event *domain.EventInfo
    result []*domain.RaceInfo
  } {
    { "no rounds", eventRoundCountsEmpty, eventInfoOneDirectFinal, raceInfosEmpty },
    // TODO add more test cases
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {
      if got, want := roundsToRaces(tt.roundCounts, tt.event), tt.result; !reflect.DeepEqual(got, want) {
        t.Errorf("raceInfo got %v, want %v", got, want)
      }
    })
  }
}

func TestRacesAndNot(t *testing.T) {
  tests := []struct{
    testName string
    r1 []*domain.RaceInfo
    r2 []*domain.RaceInfo
    result []*domain.RaceInfo
  } {
    { "r1 nil", nil, raceInfosOneDirectFinal, nil },
    { "r2 nil", raceInfosOneDirectFinal, nil, raceInfosOneDirectFinal },
    { "r1 empty", raceInfosEmpty, raceInfosOneDirectFinal, raceInfosEmpty },
    { "r2 empty", raceInfosOneDirectFinal, raceInfosEmpty, raceInfosOneDirectFinal },
    { "common", raceInfosOneDirectFinal, raceInfosOneDirectFinal, raceInfosEmpty },
    // TODO add more test cases
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {
      if got, want := racesAndNot(tt.r1, tt.r2), tt.result; !reflect.DeepEqual(got, want) {
        t.Errorf("raceInfo got %v, want %v", got, want)
      }
    })
  }
}

func TestRacesIntersectAndDiffer(t *testing.T) {
  tests := []struct{
    testName string
    r1 []*domain.RaceInfo
    r2 []*domain.RaceInfo
    result []*domain.RaceInfo
  } {
    { "r1 nil", nil, raceInfosOneDirectFinal, nil },
    { "r2 nil", raceInfosOneDirectFinal, nil, nil },
    { "r1 empty", raceInfosEmpty, raceInfosOneDirectFinal, raceInfosEmpty },
    { "r2 empty", raceInfosOneDirectFinal, raceInfosEmpty, raceInfosEmpty },
    { "common no diff", raceInfosOneDirectFinal, raceInfosOneDirectFinal, raceInfosEmpty },
    { "common diff", raceInfosOneHeat, raceInfosOneDirectFinal, raceInfosOneHeat },
    // TODO add more test cases
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {
      if got, want := racesIntersectAndDiffer(tt.r1, tt.r2), tt.result; !reflect.DeepEqual(got, want) {
        t.Errorf("raceInfo got %v, want %v", got, want)
      }
    })
  }
}

func TestFindMatchingRace(t *testing.T) {
  tests := []struct{
    testName string
    ra []*domain.RaceInfo
    r *domain.RaceInfo
    result *domain.RaceInfo
  } {
    { "find one", raceInfosOneDirectFinal, raceInfoOneDirectFinal, raceInfoOneDirectFinal },
    { "not found", raceInfosEmpty, raceInfoOneDirectFinal, nil },
    // TODO add more test cases
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {
      if got, want := findMatchingRace(tt.ra, tt.r), tt.result; !reflect.DeepEqual(got, want) {
        t.Errorf("raceInfo got %v, want %v", got, want)
      }
    })
  }
}
