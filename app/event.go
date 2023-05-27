package app

import (
  "context"
  "fmt"

  "github.com/jimmc/jraceman/domain"

  "github.com/golang/glog"
)

// CreateRacesResult is the data that we return from a CreateRaces request for an event.
type CreateRacesResult struct {
  EventRaces *domain.EventRaces
  RacesToCreate []*domain.RaceInfo
  RacesToDelete []*domain.RaceInfo
  RacesToModFrom []*domain.RaceInfo
  RacesToModTo []*domain.RaceInfo
}

// EventCreateRaces creates or updates the races for the given event
// and specified number of lanes. It may create new races, or delete or
// update existing races.
func EventCreateRaces(ctx context.Context, r domain.Repos, eventId string, laneCount int, dryRun bool, allowDeleteLanes bool) (*CreateRacesResult, error) {
  eventRaces, err := r.EventRaces().EventRaceInfo(eventId)
  if err != nil {
    return nil, fmt.Errorf("error getting races for event %q: %w", eventId, err)
  }
  glog.V(3).Infof("RaceInfo for eventId=%s: %v", eventId, eventRaces)
  // Get the progression system for the specified event.
  progression, err := ProgSysForEvent(r, eventId)
  if err != nil {
    return nil, fmt.Errorf("error getting progression system for event %q: %w", eventId, err)
  }

  result, err := calculateRaceChanges(eventRaces, progression, laneCount)
  if err != nil {
    return nil, fmt.Errorf("error calculating race changes for event %q: %w", eventId, err)
  }

  if dryRun {
    return result, nil
  }

  if !allowDeleteLanes {
    raceWithLaneData := firstRaceWithLaneData(result.RacesToDelete)
    if raceWithLaneData != nil {
      return nil, fmt.Errorf("attempt to delete a race with lane data (%s), with allowDeleteLanes false", raceWithLaneData.RaceID)
    }
  }

  // We are clear to proceed. We can now create, delete, and modify the races.
  err = r.EventRaces().UpdateRaceInfo(ctx, eventRaces, result.RacesToCreate, result.RacesToDelete,
      result.RacesToModFrom, result.RacesToModTo)
  if err!=nil {
    return nil, err
  }

  return result, nil
}

func calculateRaceChanges(eventRaces *domain.EventRaces, progression ProgSys, laneCount int) (*CreateRacesResult, error) {
  result := &CreateRacesResult{
    EventRaces: eventRaces,
  }
  if laneCount < 0 {
    laneCount = eventRaces.EntryCount
    if eventRaces.GroupSize > 1 {
      laneCount = eventRaces.GroupCount
    }
  }
  var desiredRoundCounts []*domain.EventRoundCounts
  var err error
  if laneCount == 0 {
    // If no lanes, then there will be no races.
    desiredRoundCounts = make([]*domain.EventRoundCounts,0)
  } else {
    desiredRoundCounts, err = progression.DesiredRoundCounts(
          eventRaces.ProgressionState, laneCount, eventRaces.AreaLanes, eventRaces.AreaExtraLanes)
    if err != nil {
      return nil, err
    }
  }
  existingRoundCounts := eventRaces.RoundCounts
  glog.V(3).Infof("desiredRoundCounts=%v, existingRoundCounts=%v", desiredRoundCounts, existingRoundCounts)

  // If there are no existing races, and there are no entries, that's an error.
  // The event should be scratched.
  if len(existingRoundCounts)==0 && len(desiredRoundCounts)==0 {
    return nil, fmt.Errorf("no entries and no existing races for event %s", eventRaces.Summary)
  }

  // Figure out what races we need to create, delete, or update.
  // We are asuming existingRaces is sorted by round and section.
  existingRaces := eventRaces.Races
  // We ae assuming desiredRoundCounts is sorted by round,
  // so that desiredRaces is sorted by round and section.
  desiredRaces := roundsToRaces(desiredRoundCounts, eventRaces)
  racesToCreate := racesAndNot(desiredRaces, existingRaces)
  racesToDelete := racesAndNot(existingRaces, desiredRaces)
  racesToModFrom := racesIntersectAndDiffer(existingRaces, desiredRaces)
  racesToModTo := racesIntersectAndDiffer(desiredRaces, existingRaces)
  if len(racesToModFrom) != len(racesToModTo) {
    return nil, fmt.Errorf("length mismach in races to mod lists")
      // Programming error, should not happen.
  }
  glog.V(3).Infof("racesToCreate: %v", racesToCreate)
  glog.V(3).Infof("racesToDelete: %v", racesToDelete)
  glog.V(3).Infof("racesToModFrom: %v", racesToModFrom)
  glog.V(3).Infof("racesToModTo: %v", racesToModTo)

  result.RacesToCreate = racesToCreate
  result.RacesToDelete = racesToDelete
  result.RacesToModFrom = racesToModFrom
  result.RacesToModTo = racesToModTo

  return result, nil
}

// CalculateRaceChangesForTesting is solely to allow unit testing
// of calculateRaceChanges.
func CalculateRaceChangesForTesting(eventRaces *domain.EventRaces, progression ProgSys, laneCount int) (*CreateRacesResult, error) {
  return calculateRaceChanges(eventRaces, progression, laneCount)
}

// wouldDeleteLanes returns true if any of the races have lane data.
func firstRaceWithLaneData(races []*domain.RaceInfo) *domain.RaceInfo {
  for _, race := range races {
    if race.LaneCount > 0 {
      return race
    }
  }
  return nil
}

// roundToRaces takes a list of round counts and produces a slice of
// RaceInfo structs that have the appropriate Round and Section fields filled in.
func roundsToRaces(desiredRoundCounts []*domain.EventRoundCounts, eventRaces *domain.EventRaces) []*domain.RaceInfo {
  result := make([]*domain.RaceInfo, 0)
  for _, rc := range desiredRoundCounts {
    // Round and Section are both 1-based numbers.
    for n := 1; n <= rc.Count; n++ {
      race := &domain.RaceInfo{
        EventID: eventRaces.EventID,
        AreaID: eventRaces.AreaID,
        AreaName: eventRaces.AreaName,
        StageID: rc.StageID,
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

// racesIntersectAndDiffer return a slice containing all of the RaceInfo entries in r1 that
// have a matching RaceInfo in r2 based only on the round and section fields,
// and where there is a difference in the StageNumber.
func racesIntersectAndDiffer(r1, r2 []*domain.RaceInfo) []*domain.RaceInfo {
  if r1==nil {
    return nil
  }
  if r2==nil {
    return nil
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
