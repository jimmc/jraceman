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
    return nil, fmt.Errorf("can't find event %q to get progression system: %w", eventId, err)
  }
  if event.ProgressionID == nil || *event.ProgressionID == ""{
    return nil, fmt.Errorf("event %q has no progression specified", eventId)
  }
  progressionId := *event.ProgressionID
  progression, err := r.Progression().FindByID(progressionId)
  if err != nil {
    return nil, fmt.Errorf("can't find progression %q for event %q: %w", progressionId, eventId, err)
  }
  switch progression.Class {
    case "":
      return nil, fmt.Errorf("progression %q has no class specified", progressionId)
    case "ProgressionSimplan":
      progSys, err := r.SimplanSys().LoadSimplanSys(progression)
      if err!=nil {
        return nil, fmt.Errorf("error loading simplan for progression %q: %w", progressionId, err)
      }
      return progSys, nil
    case "ProgressionComplan":
      return nil, fmt.Errorf("ProgressionForEvent: ProgressionComplan NYI")
    default:
      return nil, fmt.Errorf("progression %q has unknown progression class %q",
          progressionId, progression.Class)
  }
}
