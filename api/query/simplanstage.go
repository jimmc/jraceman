package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type simplanstageQuery struct{
  h *handler
}

func (sc *simplanstageQuery) EntityTypeName() string {
  return "simplanstage"
}

func (sc *simplanstageQuery) NewEntity() interface{} {
  return &domain.SimplanStage{}
}

func (h *handler) simplanstage(w http.ResponseWriter, r *http.Request) {
  sq := &simplanstageQuery{h}
  h.stdquery(w, r, sq)
}
