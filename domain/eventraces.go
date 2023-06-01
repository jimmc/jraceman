package domain

import (
  "context"
  "fmt"
)

// EventRacesRepo describes how various types related to event races are loaded and saved.
type EventRacesRepo interface {
  EventRaceInfo(ID string) (*EventRaces, error)
  UpdateRaceInfo(ctx context.Context, eventRaces *EventRaces, racesToCreate, racesToDelete, racesToModFrom, racesToModTo []*RaceInfo) error
}

type EventRoundCounts struct {
  Count int
  Round int
  StageID string
  StageName string
  StageNumber int
  IsFinal bool
}

func (r *EventRoundCounts) String() string {
  return fmt.Sprintf("{count=%d,round=%d,stage=%s}", r.Count, r.Round, r.StageName)
}

// EventRaces is a summary of an event with details collcted from multiple tables.
type EventRaces struct {
  EventID string
  EntryCount int
  GroupCount int
  GroupSize int
  Summary string
  RoundCounts []*EventRoundCounts
  Races []*RaceInfo
  AreaID string
  AreaName string
  AreaLanes int
  AreaExtraLanes int
  ProgressionState string
}
