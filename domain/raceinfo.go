package domain

import (
  "fmt"
)

// RaceInfo collects information about a race from multiple tables.
type RaceInfo struct {
  RaceID string
  EventID string
  StageID string
  AreaID string
  StageName string
  StageNumber int
  IsFinal bool
  Round int
  Section int
  AreaName string
  RaceNumber int
  LaneCount int         // The number of lane records associated with this race
}

func (r *RaceInfo) String() string {
  finalstr := ""
  if r.IsFinal {
    finalstr = "F"
  }
  return fmt.Sprintf("{stage=%s/%d%s,round=%d,section=%d,area=%s,racenumber=%d}",
    r.StageName, r.StageNumber, finalstr, r.Round, r.Section, r.AreaName, r.RaceNumber)
}
