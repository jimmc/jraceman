package app

import (
  "fmt"

  "github.com/jimmc/jraceman/domain"

  "github.com/golang/glog"
)

// CreateRacesResult is the data that we return from a CreateRaces request for an event.
type CreateRacesResult struct {
  EventInfo *domain.EventInfo
  RacesToCreate []*domain.RaceInfo
  RacesToDelete []*domain.RaceInfo
  RacesToModFrom []*domain.RaceInfo
  RacesToModTo []*domain.RaceInfo
}

// EventCreateRaces creates or updates the races for the given event
// and specified number of lanes. It may create new races, or delete or
// update existing races.
func EventCreateRaces(r domain.Repos, eventId string, laneCount int, dryRun bool) (*CreateRacesResult, error) {
  eventInfo, err := r.EventInfo().EventRaceInfo(eventId)
  if err != nil {
    return nil, err
  }
  result := &CreateRacesResult{
    EventInfo: eventInfo,
  }
  glog.V(3).Infof("RaceInfo for eventId=%s: %v", eventId, eventInfo)
  eventLaneCount := eventInfo.EntryCount
  if eventInfo.GroupSize > 1 {
    eventLaneCount = eventInfo.GroupCount
  }
  if laneCount < 0 {
    laneCount = eventLaneCount
  }
  var desiredRoundCounts []*domain.EventRoundCounts
  if laneCount == 0 {
    // If no lanes, then there will be no races.
    desiredRoundCounts = make([]*domain.EventRoundCounts,0)
  } else {
    // Get the progression system for the specified event.
    progression, err := ProgSysForEvent(r, eventId)
    if err != nil {
      return nil, err
    }
    desiredRoundCounts, err = progression.DesiredRoundCounts(
          eventInfo.ProgressionState, laneCount, eventInfo.AreaLanes, eventInfo.AreaExtraLanes)
    if err != nil {
      return nil, err
    }
  }
  existingRoundCounts := eventInfo.RoundCounts
  glog.V(3).Infof("desiredRoundCounts=%v, existingRoundCounts=%v", desiredRoundCounts, existingRoundCounts)

  // Figure out what races we need to create, delete, or update.
  // We are asuming existingRaces is sorted by round and section.
  existingRaces := eventInfo.Races
  // We ae assuming desiredRoundCounts is sorted by round,
  // so that desiredRaces is sorted by round and section.
  desiredRaces := roundsToRaces(desiredRoundCounts)
  racesToCreate := racesAndNot(desiredRaces, existingRaces)
  racesToDelete := racesAndNot(existingRaces, desiredRaces)
  racesToModFrom := racesIntersectAndDiffer(existingRaces, desiredRaces)
  racesToModTo := racesIntersectAndDiffer(desiredRaces, existingRaces)
  if len(racesToModFrom) != len(racesToModTo) {
    return nil, fmt.Errorf("Length mismach in races to mod lists")
      // Programming error, should not happen.
  }
  glog.V(3).Infof("racesToCreate: %v", racesToCreate)
  glog.V(3).Infof("racesToDelete: %v", racesToDelete)
  glog.V(3).Infof("racesToModFrom: %v", racesToModFrom)
  glog.V(3).Infof("racesToModTo: %v", racesToModTo)

  if dryRun {
    result.RacesToCreate = racesToCreate
    result.RacesToDelete = racesToDelete
    result.RacesToModFrom = racesToModFrom
    result.RacesToModTo = racesToModTo
    return result, nil
  }
  // TODO - check to see if the existing races have any data, return an error if
  // so unless we have a "deleteExistingLanes" flag telling us to delete that data.
  return nil, fmt.Errorf("EventCreateRaces NYI")
}

// roundToRaces takes a list of round counts and produces a slice of
// RaceInfo structs that have the appropriate Round and Section fields filled in.
func roundsToRaces(desiredRoundCounts []*domain.EventRoundCounts) []*domain.RaceInfo {
  result := make([]*domain.RaceInfo, 0)
  for _, rc := range desiredRoundCounts {
    // Round and Section are both 1-based numbers.
    for n := 1; n <= rc.Count; n++ {
      race := &domain.RaceInfo{
        StageName: rc.StageName,
        StageNumber: rc.StageNumber,
        IsFinal: rc.IsFinal,
        Round: rc.Round,
        Section: n,
      }
      result = append(result, race)
    }
  }
  return result
}

// racesAndNot returns a slice containing all of the RaceInfo entries in r1 that
// do not have a matching RaceInfo in r2 based only on the round and section fields.
func racesAndNot(r1, r2 []*domain.RaceInfo) []*domain.RaceInfo {
  if r1==nil {
    return nil
  }
  if r2==nil {
    return r1
  }
  result := make([]*domain.RaceInfo, 0)
  for _, r1r := range r1 {
    r2r := findMatchingRace(r2, r1r)
    if r2r == nil {
      result = append(result, r1r)
    }
  }
  return result
}

// racesIntersect return a slice containing all of the RaceInfo entries in r1 that
// have a matching RaceInfo in r2 based only on the round and section fields,
// and where there is a difference in the StageNumber.
func racesIntersectAndDiffer(r1, r2 []*domain.RaceInfo) []*domain.RaceInfo {
  if r1==nil {
    return nil
  }
  if r2==nil {
    return r1
  }
  result := make([]*domain.RaceInfo, 0)
  for _, r1r := range r1 {
    r2r := findMatchingRace(r2, r1r)
    if r2r != nil {
      if r2r.StageNumber != r1r.StageNumber {
        result = append(result, r1r)
      }
    }
  }
  return result
}

// findMatchingRace looks through ra for a race that matches r
// based only on the round and section fields.
func findMatchingRace(ra []*domain.RaceInfo, r *domain.RaceInfo) *domain.RaceInfo {
  for _, t := range ra {
    if t.Round == r.Round && t.Section == r.Section {
      return t
    }
  }
  return nil
}
