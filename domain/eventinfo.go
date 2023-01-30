package domain

import (
  "fmt"
)

// EvenInfoRepo describes how various types related to events are loaded and saved.
type EventInfoRepo interface {
  EventRaceInfo(ID string) (*EventInfo, error)
}

type RaceCountInfo struct {
  Count int
  Round int
  StageName string
}

func (r *RaceCountInfo) String() string {
  return fmt.Sprintf("{count=%d,round=%d,stage=%s}", r.Count, r.Round, r.StageName)
}

type EventInfo struct {
  EntryCount int
  GroupCount int
  GroupSize int
  Summary string
  RaceCounts []*RaceCountInfo
}
