package dbrepo

import (
  "fmt"
  "strings"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/domain"

  "github.com/golang/glog"
)

type DBSimplanSysRepo struct {
  db conn.DB
}

func (r *DBSimplanSysRepo) LoadSimplanSys(progression *domain.Progression) (*domain.SimplanSys, error) {
  s := &domain.SimplanSys{}
  // 1. Extract parameter values from progression.Parameters
  parameters, err := progression.ParmsAsMap()
  if err != nil {
    return nil, fmt.Errorf("SymplanSys: failed to convert progression params %q to map: %w", progression.ID, err)
  }
  system, ok := parameters["system"]
  if !ok || system=="" {
    return nil, fmt.Errorf("SimplanSys: no system name in parameters for progressionId %q", progression.ID)
  }
  s.System = system
  multipleFinals, ok := parameters["multipleFinals"]
  if ok {
    b, ok := parameterToBoolean(multipleFinals)
    if !ok {
      return nil, fmt.Errorf("invalid value %q for mulitpleFinals option in progression %q",
          multipleFinals, progression.ID)
    }
    s.MultipleFinals = b
  }
  useExtraLanes, ok := parameters["useExtraLanes"]
  if ok {
    b, ok := parameterToBoolean(useExtraLanes)
    if !ok {
      return nil, fmt.Errorf("invalid value %q for useExtraLanes option in progression %q",
          useExtraLanes, progression.ID)
    }
    s.UseExtraLanes = b
  }

    _, err = r.db.Query("SELECT COUNT(*) FROM SimplanStage")
    if err != nil {
      return nil, fmt.Errorf("failed to count rows in SimplanStage: %w", err)
    }

  // 2. Get plan info for all plans for this system
  query := `SELECT ID, Plan, MinEntries, maxEntries from Simplan
      where System=?`
  whereVals := make([]interface{}, 1)
  whereVals[0] = system
  glog.V(3).Infof("SQL: %s with whereVals=%#v", query, whereVals)
  rows, err := r.db.Query(query, whereVals...)
  if err != nil {
    return nil, fmt.Errorf("sql query to read Simplan for system %q failed in progression %q: %w",
        system, progression.ID, err)
  }
  defer rows.Close()
  plans := make([]*domain.SimplanRoundsPlan,0)
  for rows.Next() {
    plan := &domain.SimplanRoundsPlan{}
    err := rows.Scan(&plan.SimplanID, &plan.Name, &plan.MinEntries, &plan.MaxEntries)
    if err != nil {
      return nil, fmt.Errorf("row scan failed for Simplan system %q in progression %q: %w",
        system, progression.ID, err)
    }

    // 3. Get the stageid and sectioncount from all rows in the simplanstage
    //    table with the simplan id from the previous step
    stagesQuery := `SELECT Stage.ID as StageId, SimplanStage.SectionCount as SectionCount,
              Stage.Name as StageName, Stage.Number as StageNumber, Stage.IsFinal as IsFinal
            FROM SimplanStage JOIN Stage on SimplanStage.StageID=Stage.ID
            WHERE SimplanStage.SimplanID=?
            ORDER BY StageNumber`
    stagesVals := make([]interface{}, 1)
    stagesVals[0] = plan.SimplanID
    glog.V(3).Infof("SQL: %s with whereVals=%#v", stagesQuery, stagesVals)

    stageRows, err := r.db.Query(stagesQuery, stagesVals...)
    if err != nil {
      return nil, fmt.Errorf("failed to read SimplanStage for simplan %q in progression %q: %w",
        plan.SimplanID, progression.ID, err)
    }
    defer stageRows.Close()
    round := 1    // Round is 1-based.
    rowCount := 0
    raceCounts := make([]*domain.EventRoundCounts,0)
    for stageRows.Next() {
      rci := &domain.EventRoundCounts{}
      err := stageRows.Scan(&rci.StageID, &rci.Count, &rci.StageName, &rci.StageNumber, &rci.IsFinal)
      if err != nil {
        return nil, fmt.Errorf("failed to read a SimplanStage row for simplan %q in progression %q: %w",
            plan.SimplanID, progression.ID, err)
      }
      rci.Round = round
      round++
      raceCounts = append(raceCounts, rci)
      rowCount++
    }
    plan.RoundCounts = raceCounts
    plans = append(plans, plan)
  }
  s.Rounds = plans

  return s, nil
}

// Parses a string as a boolean value. On error, the second return value is false.
func parameterToBoolean(v string) (bool, bool) {
  switch strings.ToLower(v) {
    case "true", "on", "yes", "": return true, true
    case "false", "off", "no": return false, true
    default: return false, false
  }
}
