package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
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

func (h *handler) laneorder(w http.ResponseWriter, r *http.Request) {
  sq := &laneorderQuery{h}
  h.stdquery(w, r, sq)
}
