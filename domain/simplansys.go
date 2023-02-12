package domain

import (
  "fmt"
)

// SimplanSysRepo describes how composite simplan structs are loaded and saved.
type SimplanSysRepo interface {
  LoadSimplanSys(progression *Progression) (*SimplanSys, error)
}

// SimplanRoundsPlan implements the ProgRoundsPlan interface for one plan
// according to the rules for a simple progression plan.
type SimplanRoundsPlan struct {
  System SimplanSys     // The system containing this plan
  SimplanID string      // The ID of this plan
  Name string           // Plan name
  MinEntries int        // Minimum number of entries to select this plan
  MaxEntries int        // Maximum number 0f entireto select this plan
  RoundCounts []*EventRoundCounts  // The count of the number of races we should have per round
}

// SimplanSys implements the ProgSys interface according to the
// rules for a simple progression plan.
type SimplanSys struct {
  System string         // The system name for this progression plan
  MultipleFinals bool   // True means use multiple direct finals instead of heats
  UseExtraLanes bool    // True means use extra lanes when an area has them
  Rounds []*SimplanRoundsPlan   // The possible plans for this system
}

func (s *SimplanSys) DesiredRoundCounts(progressionState string, laneCount, areaLanes, areaExtraLanes int) ([]*EventRoundCounts, error) {
  if laneCount <= areaLanes || (s.UseExtraLanes && laneCount <= (areaLanes + areaExtraLanes)) {
    // Everything fits into one direct final. Find the "DF" plan.
    for _, p := range s.Rounds {
      if p.Name == "DF" {
        return p.RoundCounts, nil
      }
    }
    return nil, fmt.Errorf("Could not find plan DF for system %q", s.System)
  }
  if s.MultipleFinals {
    // Use however many direct finals it takes.
    numFinals := (laneCount + (areaLanes - 1)) / areaLanes
    // Find the direct final plan.
    var p0 SimplanRoundsPlan
    for _, p := range s.Rounds {
      if p.Name == "DF" {
        p0 = *p
        break
      }
    }
    if p0.Name == "" {
      return nil, fmt.Errorf("Could not find plan DF for system %q", s.System)
    }
    if len(p0.RoundCounts) != 1 {
      return nil, fmt.Errorf("Expected DF plan %q to have 1 round, but it has %d", p0.SimplanID, len(p0.RoundCounts))
    }
    roundCounts := *p0.RoundCounts[0]  // Make a copy we can modify.
    roundCounts.Count = numFinals
    a := make([]*EventRoundCounts,1)
    a[0] = &roundCounts
    return a, nil
  }
  // TODO - if a plan is specified in progressionState, and it has the appropriate
  // lane count, use it.
  // Look for a plan with the right lane count.
  for _, p := range s.Rounds {
    if laneCount <= p.MaxEntries && laneCount >= p.MinEntries {
      return p.RoundCounts, nil
    }
  }
  return nil, fmt.Errorf("Could not find plan DF for system %q", s.System)
}
