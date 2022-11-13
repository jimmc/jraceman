package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type raceQuery struct{
  domain.RaceMeta
  h *handler
}

func (sc *raceQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) race(w http.ResponseWriter, r *http.Request) {
  sq := &raceQuery{domain.RaceMeta{}, h}
  h.stdquery(w, r, sq)
}
