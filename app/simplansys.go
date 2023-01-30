package app

import (
  "database/sql"
  "fmt"

  "github.com/jimmc/jraceman/dbrepo"
  "github.com/jimmc/jraceman/domain"

  "github.com/golang/glog"
)

// SimplanSys implements the ProgSys interface according to the
// rules for a simple progression plan.
type SimplanSys struct {
  system string         // The system name for this progression plan
  simplanID string      // The ID of the Simplan entry we are using
  laneCount int         // The number of lanes to use when creating races
  raceCounts []*RaceCountInfo  // The count of the number of races we should have per round
}

func NewSimplanSys(dbr *dbrepo.Repos, progression *domain.Progression, progressionState *string, laneCount int) (*SimplanSys, error) {
  s := &SimplanSys{}
  // 1. Extract system name from progression.Parameters
  parameters, err := progressionParmsToMap(progression.Parameters)
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
  err = dbr.DB().QueryRow(query, whereVals...).Scan(&s.simplanID)
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
  stagesVals[0] = s.simplanID
  glog.V(3).Infof("SQL: %s with whereVals=%#v", stagesQuery, stagesVals)

  rows, err := dbr.DB().Query(stagesQuery, stagesVals...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  rowCount := 0
  raceCounts := make([]*RaceCountInfo,0)
  for rows.Next() {
    stageId := ""
    stageNumber := 0
    rci := &RaceCountInfo{}
    err := rows.Scan(&stageId, &rci.Count, &rci.StageName, &stageNumber)
    if err != nil {
      return nil, err
    }
    raceCounts = append(raceCounts, rci)
    rowCount++
  }
  s.raceCounts = raceCounts

  return s, nil
}

func (s *SimplanSys) DesiredRaceCounts() ([]*RaceCountInfo, error) {
  return s.raceCounts, nil
}
