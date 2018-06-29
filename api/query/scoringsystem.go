package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type scoringsystemQuery struct{
  h *handler
}

func (sc *scoringsystemQuery) EntityTypeName() string {
  return "scoringsystem"
}

func (sc *scoringsystemQuery) NewEntity() interface{} {
  return &domain.ScoringSystem{}
}

func (h *handler) scoringsystem(w http.ResponseWriter, r *http.Request) {
  sq := &scoringsystemQuery{h}
  h.stdquery(w, r, sq)
}
