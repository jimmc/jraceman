package domain

// SimplanSysRepo describes how composite simplan structs are loaded and saved.
type SimplanSysRepo interface {
  LoadSimplanSys(progression *Progression, progressionState *string, laneCount int) (*SimplanSys, error)
}

// SimplanSys implements the ProgSys interface according to the
// rules for a simple progression plan.
type SimplanSys struct {
  System string         // The system name for this progression plan
  SimplanID string      // The ID of the Simplan entry we are using
  LaneCount int         // The number of lanes to use when creating races
  RaceCounts []*RaceCountInfo  // The count of the number of races we should have per round
}

func (s *SimplanSys) DesiredRaceCounts() ([]*RaceCountInfo, error) {
  return s.RaceCounts, nil
}
