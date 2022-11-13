package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type laneQuery struct{
  domain.LaneMeta
  h *handler
}

func (sc *laneQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) lane(w http.ResponseWriter, r *http.Request) {
  sq := &laneQuery{domain.LaneMeta{}, h}
  h.stdquery(w, r, sq)
}
