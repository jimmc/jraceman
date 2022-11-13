package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type eventQuery struct{
  domain.EventMeta
  h *handler
}

func (sc *eventQuery) SummaryQuery(format string) string {
  switch format {
  case "byname":
    return "select ID, Name || ' [' || ID || ']' as summary from (select * from event where number>0) ORDER BY Name"
  case "bynumber":
    return "select ID, '#' || number || ' ' || Name || ' [' || ID || ']' as summary from (select * from event where number>0) ORDER BY number"
  case "byid":
    return "select ID, '[' || ID || '] ' || Name as summary from (select * from event where number>0) ORDER BY ID"
  case "byracenumber":
    return `select event.ID,
              'Race #' || race.number || ': ' || event.Name || ' [' || event.ID || '] ' || stage.name || ' ' || race.round as summary
            from (select * from event where number>0) as event JOIN race on event.id = race.eventid LEFT JOIN stage on race.stageid = stage.id
            ORDER BY race.number`
  default: // Includes unknown formats and default (no format).
    return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
  }
}

func (h *handler) event(w http.ResponseWriter, r *http.Request) {
  sq := &eventQuery{domain.EventMeta{}, h}
  h.stdquery(w, r, sq)
}
