package dbrepo

import (
  "context"
  "database/sql"
  "fmt"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/domain"

  "github.com/golang/glog"
)

type DBEventInfoRepo struct {
  db conn.DB
  repos *Repos      // We need access to repos for the table types
}

func (r *DBEventInfoRepo) EventRaceInfo(eventId string) (*domain.EventInfo, error) {
  if eventId == "" {
    return nil, fmt.Errorf("Event ID must be specified")
  }
  db := r.db

  // Collect summary info about the entries for this event.
  entryCountQuery :=
        `(SELECT count(1) as entryCount
        FROM entry JOIN event on entry.eventid = event.id
        WHERE event.id=? AND NOT entry.scratched)`
  groupCountQuery :=
        `(SELECT count(distinct groupname) as groupCount
        FROM entry JOIN event on entry.eventid = event.id
        WHERE event.id=? AND NOT entry.scratched)`
  query := "SELECT "+entryCountQuery+" as EntryCount,"+
        groupCountQuery+` as GroupCount,
        COALESCE(competition.groupsize,0) as GroupSize,
        event.Name || ' [' || event.ID || ']' as Summary,
        COALESCE(event.areaid,"") as areaID,
        COALESCE(area.Name,"") as areaName,
        COALESCE(area.Lanes,0) as areaLanes,
        COALESCE(area.ExtraLanes,0) as areaExtraLanes,
        COALESCE(event.ProgressionState,"") as progressionState
        FROM event
        LEFT JOIN competition on event.competitionid = competition.id
        LEFT JOIN area on event.areaid = area.id
        WHERE event.id=?`
  whereVals := make([]interface{}, 3)
  whereVals[0] = eventId
  whereVals[1] = eventId
  whereVals[2] = eventId
  glog.V(3).Infof("SQL: %s", query)
  result := &domain.EventInfo{}
  err := db.QueryRow(query, whereVals...).Scan(
    &result.EntryCount, &result.GroupCount, &result.GroupSize, &result.Summary,
    &result.AreaID, &result.AreaName, &result.AreaLanes, &result.AreaExtraLanes,
    &result.ProgressionState)
  if err != nil {
    return nil, fmt.Errorf("Error collecting event %q info: %w", eventId, err)
  }
  result.EventID = eventId

  // Collect information about the races for this event.
  races, err := loadEventRaces(db, eventId)
  if err != nil {
    return nil, fmt.Errorf("Error collecting races for event %q: %w", eventId, err)
  }
  result.Races = races

  // Get the count of races that exist for this event.
  roundCounts, err := loadEventRoundCounts(db, eventId)
  if err != nil {
    return nil, fmt.Errorf("Error collecting round counts for event %q: %w", eventId, err)
  }
  result.RoundCounts = roundCounts
  return result, nil
}

func loadEventRaces(db conn.DB, eventId string) ([]*domain.RaceInfo, error) {
  laneCountQuery :=
        `(SELECT lCount FROM
          (SELECT count(1) as lCount, lrace.id as lraceid
          FROM lane JOIN race as lrace on lane.raceid = lrace.id
          GROUP BY lrace.id) as LaneCounts
        WHERE LaneCounts.lraceid = race.id)`
  query := `SELECT COALESCE(stage.name,"") as StageName,
        COALESCE(stage.number,0) as StageNumber, COALESCE(stage.isfinal,false) as IsFinal,
        race.round as Round, race.section as Section,
        COALESCE(area.name,"") as AreaName, COALESCE(race.number,0) as RaceNumber, race.ID as RaceID,
        COALESCE(`+laneCountQuery+`,0) as LaneCount,
        COALESCE(race.stageid,"") as StageID, COALESCE(race.areaid,"") as AreaID
    FROM race LEFT JOIN stage on race.stageid = stage.id
        LEFT JOIN area on race.areaid = area.id
    WHERE race.eventid = ?
    ORDER BY race.round, race.section`
  whereVals := make([]interface{}, 1)
  whereVals[0] = eventId
  glog.V(3).Infof("SQL: %s; with whereVals=%v", query, whereVals)
  rows, err := db.Query(query, whereVals...)
  if err != nil {
    return nil, fmt.Errorf("Error collecting races for event %q: %w", eventId, err)
  }
  defer rows.Close()
  rr := make([]*domain.RaceInfo,0)
  for rows.Next() {
    r := &domain.RaceInfo{}
    if err = rows.Scan(&r.StageName, &r.StageNumber, &r.IsFinal,
        &r.Round, &r.Section, &r.AreaName, &r.RaceNumber, &r.RaceID, &r.LaneCount, &r.StageID,
        &r.AreaID); err != nil {
      return nil, fmt.Errorf("Error collecting race data for event %q: %w", eventId, err)
    }
    r.EventID = eventId
    rr = append(rr, r)
  }
  return rr, nil
}

func loadEventRoundCounts(db conn.DB, eventId string) ([]*domain.EventRoundCounts, error) {
  query := `SELECT count(1) as count, race.round as round,
      race.stageid as StageID, stage.name as StageName
    FROM event JOIN race on event.id = race.eventid
      JOIN stage on race.stageid=stage.id
    WHERE event.id=?
    GROUP BY race.round
    ORDER BY race.round
    `
  whereVals := make([]interface{}, 1)
  whereVals[0] = eventId
  glog.V(3).Infof("SQL: %s; with whereVals=%v", query, whereVals)
  rows, err := db.Query(query, whereVals...)
  if err != nil {
    return nil, fmt.Errorf("Error collecting race info for event %q: %w", eventId, err)
  }
  defer rows.Close()
  rr := make([]*domain.EventRoundCounts,0)
  for rows.Next() {
    r := &domain.EventRoundCounts{}
    if err = rows.Scan(&r.Count, &r.Round, &r.StageID, &r.StageName); err != nil {
        return nil, fmt.Errorf("Error collecting race count row for event %q: %w", eventId, err)
    }
    rr = append(rr, r)
  }
  return rr, nil
}

// UpdateRaceInfo updates the database to create, delete, and modify races
// according to the given data.
func (r *DBEventInfoRepo) UpdateRaceInfo(ctx context.Context, eventInfo *domain.EventInfo,
    racesToCreate, racesToDelete, racesToModFrom, racesToModTo []*domain.RaceInfo) error {
  // See if our database connection is already a transaction.
  tx, ok := r.db.(*sql.Tx)
  if ok {
    return r.updateRaceInfoInTx(ctx, tx, eventInfo, racesToCreate, racesToDelete, racesToModFrom, racesToModTo)
  }
  // We do all operations within a transaction and roll back if any fail.
  db, ok := r.db.(*sql.DB)
  if !ok {
    return fmt.Errorf("In EventInfo.UpdateRaceInfo db is not a Tx or DB!")
  }
  tx, err := db.BeginTx(ctx, nil)
  if err!=nil {
    return err
  }
  defer tx.Rollback()   // Roll back if anything fails.

  err = r.updateRaceInfoInTx(ctx, tx, eventInfo, racesToCreate, racesToDelete, racesToModFrom, racesToModTo)
  if err != nil {
    return err
  }

  if err = tx.Commit(); err!=nil {
    return err
  }
  return nil
}

// UpdateRaceInfo updates the database to create, delete, and modify races
// according to the given data. Return error if any operations fail.
func (r *DBEventInfoRepo) updateRaceInfoInTx(ctx context.Context, tx *sql.Tx, eventInfo *domain.EventInfo,
    racesToCreate, racesToDelete, racesToModFrom, racesToModTo []*domain.RaceInfo) error {

  // Create
  for _, raceInfo := range racesToCreate {
    race := raceInfoToRace(raceInfo)
    id, err := r.repos.Race().Save(race)
    if err!=nil {
      return fmt.Errorf("error creating race for event ID=%q: %w", raceInfo.EventID, err)
    }
    raceInfo.RaceID = id
  }

  // Delete
  for _, raceInfo := range racesToDelete {
    if err := r.repos.Race().DeleteByID(raceInfo.RaceID); err!=nil {
      return fmt.Errorf("error deleting race ID=%q: %w", raceInfo.RaceID, err)
    }
  }

  // Update
  for i, raceInfo := range racesToModFrom {
    oldRace := raceInfoToRace(raceInfo)
    newRace := raceInfoToRace(racesToModTo[i])
    diffMap := make(map[string]interface{})
    diffMap["StageID"] = newRace.StageID
    // TODO - need to deal with Race.Scratched?
    raceDiffs := &manualDiffs{ diffMap }
    err := r.repos.Race().UpdateByID(oldRace.ID, oldRace, newRace, raceDiffs)
    if err!=nil {
      return fmt.Errorf("error updating race ID=%q: %w", raceInfo.RaceID, err)
    }
    raceInfo.RaceID = oldRace.ID
  }

  return nil
}

type manualDiffs struct {
  diffMap map[string]interface{}
}

func (md *manualDiffs) Modified() map[string]interface{} { return md.diffMap }

// Create a Race struct that we can use with the database.
func raceInfoToRace(raceInfo *domain.RaceInfo) *domain.Race {
  race := &domain.Race{
    ID: raceInfo.RaceID,
    EventID: raceInfo.EventID,
    Round: &raceInfo.Round,     // We know we have Round and Section.
    Section: &raceInfo.Section,
  }
  if raceInfo.AreaID != "" {
    race.AreaID = &raceInfo.AreaID
  }
  if raceInfo.RaceNumber != 0 {
    race.Number = &raceInfo.RaceNumber
  }
  if raceInfo.StageID != "" {
    race.StageID = &raceInfo.StageID
  }
  return race
}
