package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type laneorderQuery struct{
  h *handler
}

func (sc *laneorderQuery) EntityTypeName() string {
  return "laneorder"
}

func (sc *laneorderQuery) NewEntity() interface{} {
  return &domain.LaneOrder{}
}

func (sc *laneorderQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) laneorder(w http.ResponseWriter, r *http.Request) {
  sq := &laneorderQuery{h}
  h.stdquery(w, r, sq)
}
