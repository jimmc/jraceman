package app

import (
  "fmt"

  "github.com/jimmc/jraceman/domain"
)

// ProgSys declares the methods used by different implementations
// of progression systems. A ProgSys is responsible for determining how
// three steps are implemented:
//  1. Determining the count of races in each round when creating races for an event
//  2. Seeding the lane assignments for the first round of an event
//  3. Defining how winners from one round are assigned lanes (progressed) in the next round
type ProgSys interface {
  DesiredRoundCounts(progressionState string, laneCount, areaLanes, areaExtraLanes int) ([]*domain.EventRoundCounts, error)
}

// ProgSysForEvent looks up the progression in the database for the specified event
// and returns the right kind of ProgSys for it.
func ProgSysForEvent(r domain.Repos, eventId string) (ProgSys, error) {
  event, err := r.Event().FindByID(eventId)
  if err != nil {
    return nil, err
  }
  if event.ProgressionID == nil || *event.ProgressionID == ""{
    return nil, fmt.Errorf("EventID %q has no progression specified", eventId)
  }
  progressionId := *event.ProgressionID
  progression, err := r.Progression().FindByID(progressionId)
  if err != nil {
    return nil, err
  }
  switch progression.Class {
    case "":
      return nil, fmt.Errorf("ProgressionID %q has no class specified", progressionId)
    case "ProgressionSimplan":
      return r.SimplanSys().LoadSimplanSys(progression)
    case "ProgressionComplan":
      return nil, fmt.Errorf("ProgressionForEvent: ProgressionComplan NYI")
    default:
      return nil, fmt.Errorf("ProgressionID %q has unknown progression class %q",
          progressionId, progression.Class)
  }
}
