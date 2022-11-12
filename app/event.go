package app

import (
  "fmt"

  "github.com/jimmc/jraceman/dbrepo"

  "github.com/golang/glog"
)

type RaceCountInfo struct {
  Count int
  Round int
  StageName string
}

type EventInfo struct {
  EntryCount int
  GroupCount int
  GroupSize int
  Summary string
  RaceCounts []*RaceCountInfo
}

func EventRaceInfo(dbr *dbrepo.Repos, eventId string) (*EventInfo, error) {
  if eventId == "" {
    return nil, fmt.Errorf("Event ID must be specified")
  }
  db := dbr.DB()
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
        event.Name || ' [' || event.ID || ']' as Summary
        FROM event
        LEFT JOIN competition on event.competitionid = competition.id
        WHERE event.id=?`
  whereVals := make([]interface{}, 3)
  whereVals[0] = eventId
  whereVals[1] = eventId
  whereVals[2] = eventId
  glog.V(3).Infof("SQL: %s", query)
  result := &EventInfo{}
  err := db.QueryRow(query, whereVals...).Scan(
    &result.EntryCount, &result.GroupCount, &result.GroupSize, &result.Summary)
  if err != nil {
    return nil, fmt.Errorf("Error collecting event %q info: %w", eventId, err)
  }
  // Now get the count of races that exist for this event.
  query = `SELECT count(1) as count, race.round as round, stage.name as stagename
    FROM event JOIN race on event.id = race.eventid
      JOIN stage on race.stageid=stage.id
    WHERE event.id=?
    GROUP BY race.round
    ORDER BY race.round
    `
  whereVals = make([]interface{}, 1)
  whereVals[0] = eventId
  glog.V(3).Infof("SQL: %s", query)
  rows, err := db.Query(query, whereVals...)
  if err != nil {
    return nil, fmt.Errorf("Error collecting race info for event %q: %w", eventId, err)
  }
  defer rows.Close()
  rr := make([]*RaceCountInfo,0)
  for rows.Next() {
    r := &RaceCountInfo{}
    if err = rows.Scan(&r.Count, &r.Round, &r.StageName); err != nil {
        return nil, fmt.Errorf("Error collecting race count row for event %q: %w", eventId, err)
    }
    rr = append(rr, r)
  }
  result.RaceCounts = rr
  return result, nil
}
