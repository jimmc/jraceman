package app

import (
  "fmt"
  "strings"

  "github.com/jimmc/jraceman/dbrepo"
)

// ProgSys declares the methods used by different implementations
// of progression systems. A ProgSys is responsible for determining how
// three steps are implemented:
//  1. Determining the count of races in each round when creating races for an event
//  2. Seeding the lane assignments for the first round of an event
//  3. Defining how winners from one round are assigned lanes (progressed) in the next round
type ProgSys interface {
  DesiredRaceCounts() ([]*RaceCountInfo, error)
}

// ProgSysForEvent looks up the progression in the database for the specified event
// and returns the right kind of ProgSys for it.
func ProgSysForEvent(dbr *dbrepo.Repos, eventId string, laneCount int) (ProgSys, error) {
  event, err := dbr.Event().FindByID(eventId)
  if err != nil {
    return nil, err
  }
  if event.ProgressionID == nil || *event.ProgressionID == ""{
    return nil, fmt.Errorf("EventID %q has no progression specified", eventId)
  }
  progressionId := *event.ProgressionID
  progression, err := dbr.Progression().FindByID(progressionId)
  if err != nil {
    return nil, err
  }
  switch progression.Class {
    case "":
      return nil, fmt.Errorf("ProgressionID %q has no class specified", progressionId)
    case "ProgressionSimplan":
      return NewSimplanSys(dbr, progression, event.ProgressionState, laneCount)
    case "ProgressionComplan":
      return nil, fmt.Errorf("ProgressionForEvent: ProgressionComplan NYI")
    default:
      return nil, fmt.Errorf("ProgressionID %q has unknown progression class %q",
          progressionId, progression.Class)
  }
}

// progressionParmsToMap parses the parameters string from the progression table
// and returns a map[string]string with all of the values.
// Each parameter is name=value, and parameters are separated by commas.
// There is no additional whitespace around either the equals or the commas.
func progressionParmsToMap(parmstr *string) (map[string]string, error) {
  values := make(map[string]string)
  if parmstr == nil {
    return values, nil
  }
  pkvs := strings.Split(*parmstr, ",")
  for _, kvs := range pkvs {
    kva := strings.Split(kvs, "=")
    if len(kva) != 2 {
      return nil, fmt.Errorf("Invalid syntax for progression parameter %s, should be name=value", kvs)
    } else {
      values[kva[0]] = kva[1]
    }
  }
  return values, nil
}
