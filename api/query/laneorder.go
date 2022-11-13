package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type laneorderQuery struct{
  domain.LaneOrderMeta
  h *handler
}

func (sc *laneorderQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) laneorder(w http.ResponseWriter, r *http.Request) {
  sq := &laneorderQuery{domain.LaneOrderMeta{}, h}
  h.stdquery(w, r, sq)
}
