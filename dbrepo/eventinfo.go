package dbrepo

import (
  "context"
  "database/sql"
  "fmt"

  "github.com/jimmc/jraceman/domain"

  "github.com/golang/glog"
)

type DBEventInfoRepo struct {
  db *sql.DB
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
        area.Name as areaName,
        area.Lanes as areaLanes,
        area.ExtraLanes as areaExtraLanes,
        event.ProgressionState as progressionState
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
    &result.AreaName, &result.AreaLanes, &result.AreaExtraLanes, &result.ProgressionState)
  if err != nil {
    return nil, fmt.Errorf("Error collecting event %q info: %w", eventId, err)
  }

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

func loadEventRaces(db *sql.DB, eventId string) ([]*domain.RaceInfo, error) {
  laneCountQuery :=
        `(SELECT count(1) as laneCount
        FROM lane JOIN race on lane.raceid = race.id)`
  query := `SELECT stage.name as StageName, stage.number as StageNumber, stage.isfinal as IsFinal,
        race.round as Round,
        race.section as Section, area.name as AreaName, race.number as RaceNumber, race.ID as RaceID,
        `+laneCountQuery+` as LaneCount
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
        &r.Round, &r.Section, &r.AreaName, &r.RaceNumber, &r.RaceID, &r.LaneCount); err != nil {
      return nil, fmt.Errorf("Error collecting race data for event %q: %w", eventId, err)
    }
    rr = append(rr, r)
  }
  return rr, nil
}

func loadEventRoundCounts(db *sql.DB, eventId string) ([]*domain.EventRoundCounts, error) {
  query := `SELECT count(1) as count, race.round as round, stage.name as stagename
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
    if err = rows.Scan(&r.Count, &r.Round, &r.StageName); err != nil {
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
  // We do all operations within a transaction and roll back if any fail.
  tx, err := r.db.BeginTx(ctx, nil)
  if err!=nil {
    return err
  }
  defer tx.Rollback()   // Roll back if anything fails.

  // TODO: create, delete, and update races; return error if anything fails

  if err = tx.Commit(); err!=nil {
    return err
  }
  return fmt.Errorf("EventInfo.UpdateRaceInfo NYI")
  // return nil
}
