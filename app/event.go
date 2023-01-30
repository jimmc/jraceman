package app

import (
  "fmt"

  "github.com/jimmc/jraceman/domain"

  "github.com/golang/glog"
)

// TODO - fill out CreateRacesResults
type CreateRacesResult struct {}

// EventCreateRaces creates or updates the races for the given event
// and specified number of lanes. It may create new races, or delete or
// update existing races.
func EventCreateRaces(r domain.Repos, eventId string, laneCount int) (*CreateRacesResult, error) {
  raceInfo, err := r.EventInfo().EventRaceInfo(eventId)
  if err != nil {
    return nil, err
  }
  eventLaneCount := raceInfo.EntryCount
  if raceInfo.GroupSize > 1 {
    eventLaneCount = raceInfo.GroupCount
  }
  if laneCount < 0 {
    laneCount = eventLaneCount
  }
  // Get the progression for the specified event. This will include
  // any progression state information from the event, and the number
  // of lanes required by the event.
  progression, err := ProgSysForEvent(r, eventId, laneCount)
  if err != nil {
    return nil, err
  }
  desiredRoundCounts, err := progression.DesiredRoundCounts()
  if err != nil {
    return nil, err
  }
  existingRoundCounts := raceInfo.RoundCounts
  // TODO: get existing races, get differences, calculate create/delete/update.
  glog.V(3).Infof("desiredRoundCounts=%v, existingRoundCounts=%v", desiredRoundCounts, existingRoundCounts)
  return nil, fmt.Errorf("EventCreateRaces NYI")
}
