package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type laneQuery struct{
  h *handler
}

func (sc *laneQuery) EntityTypeName() string {
  return "lane"
}

func (sc *laneQuery) NewEntity() interface{} {
  return &domain.Lane{}
}

func (sc *laneQuery) SummaryQuery() string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) lane(w http.ResponseWriter, r *http.Request) {
  sq := &laneQuery{h}
  h.stdquery(w, r, sq)
}
