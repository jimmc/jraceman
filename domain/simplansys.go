package domain

// SimplanSysRepo describes how composite simplan structs are loaded and saved.
type SimplanSysRepo interface {
  LoadSimplanSys(progression *Progression, progressionState *string, laneCount int) (*SimplanSys, error)
}

// SimplanSys implements the ProgSys interface according to the
// rules for a simple progression plan.
type SimplanSys struct {
  System string         // The system name for this progression plan
  MultipleFinals bool   // True means use multiple direct finals instead of heats
  UseExtraLanes bool    // True means use extra lanes when an area has them
  SimplanID string      // The ID of the Simplan entry we are using
  LaneCount int         // The number of lanes to use when creating races
  RoundCounts []*EventRoundCounts  // The count of the number of races we should have per round
}

func (s *SimplanSys) DesiredRoundCounts() ([]*EventRoundCounts, error) {
  return s.RoundCounts, nil
}
