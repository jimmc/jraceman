package dbrepo

import (
  "database/sql"
  "fmt"

  "github.com/jimmc/jraceman/domain"

  "github.com/golang/glog"
)

type DBSimplanSysRepo struct {
  db *sql.DB
}

func (r *DBSimplanSysRepo) LoadSimplanSys(progression *domain.Progression, progressionState *string, laneCount int) (*domain.SimplanSys, error) {
  s := &domain.SimplanSys{}
  // 1. Extract system name from progression.Parameters
  parameters, err := progression.ParmsAsMap()
  if err != nil {
    return nil, err
  }
  system, ok := parameters["system"]
  if !ok || system=="" {
    return nil, fmt.Errorf("SimplanSys: No system name in parameters for progressionId %q", progression.ID)
  }
  // 2. Get the simplan id from the row in the simplan table
  //    with matching system and with minentries<=laneCount<=maxentries
  query := `SELECT ID from Simplan
      where System=? and MinEntries<=? and MaxEntries>=?`
  whereVals := make([]interface{}, 3)
  whereVals[0] = system
  whereVals[1] = laneCount
  whereVals[2] = laneCount
  glog.V(3).Infof("SQL: %s with whereVals=%#v", query, whereVals)
  err = r.db.QueryRow(query, whereVals...).Scan(&s.SimplanID)
  if err!=nil {
    if err == sql.ErrNoRows {
      return nil, fmt.Errorf("No Simplan found for system=%q and entries=%d", system, laneCount)
    }
    return nil, err
  }

  // 3. Get the stageid and sectioncount from all rows in the simplanstage
  //    table with the simplan id from the previous step
  stagesQuery := `SELECT Stage.ID as StageId, SimplanStage.SectionCount as SectionCount,
            Stage.Name as StageName, Stage.Number as StageNumber
          FROM SimplanStage JOIN Stage on SimplanStage.StageID=Stage.ID
          WHERE SimplanStage.SimplanID=?
          ORDER BY StageNumber`
  stagesVals := make([]interface{}, 1)
  stagesVals[0] = s.SimplanID
  glog.V(3).Infof("SQL: %s with whereVals=%#v", stagesQuery, stagesVals)

  rows, err := r.db.Query(stagesQuery, stagesVals...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  rowCount := 0
  raceCounts := make([]*domain.EventRoundCounts,0)
  for rows.Next() {
    stageId := ""
    stageNumber := 0
    rci := &domain.EventRoundCounts{}
    err := rows.Scan(&stageId, &rci.Count, &rci.StageName, &stageNumber)
    if err != nil {
      return nil, err
    }
    raceCounts = append(raceCounts, rci)
    rowCount++
  }
  s.RoundCounts = raceCounts

  return s, nil
}
