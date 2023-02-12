package domain

import (
  "context"
  "fmt"
)

// EvenInfoRepo describes how various types related to events are loaded and saved.
type EventInfoRepo interface {
  EventRaceInfo(ID string) (*EventInfo, error)
  UpdateRaceInfo(ctx context.Context, eventInfo *EventInfo, racesToCreate, racesToDelete, racesToModFrom, racesToModTo []*RaceInfo) error
}

type EventRoundCounts struct {
  Count int
  Round int
  StageName string
  StageNumber int
  IsFinal bool
}

func (r *EventRoundCounts) String() string {
  return fmt.Sprintf("{count=%d,round=%d,stage=%s}", r.Count, r.Round, r.StageName)
}

// EventInfo is a summary of an event with details collcted from multiple tables.
type EventInfo struct {
  EntryCount int
  GroupCount int
  GroupSize int
  Summary string
  RoundCounts []*EventRoundCounts
  Races []*RaceInfo
  AreaName string
  AreaLanes int
  AreaExtraLanes int
  ProgressionState string
}
