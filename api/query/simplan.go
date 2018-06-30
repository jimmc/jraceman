package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type simplanQuery struct{
  h *handler
}

func (sc *simplanQuery) EntityTypeName() string {
  return "simplan"
}

func (sc *simplanQuery) NewEntity() interface{} {
  return &domain.Simplan{}
}

func (h *handler) simplan(w http.ResponseWriter, r *http.Request) {
  sq := &simplanQuery{h}
  h.stdquery(w, r, sq)
}
