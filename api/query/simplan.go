package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
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

func (sc *simplanQuery) SummaryQuery() string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) simplan(w http.ResponseWriter, r *http.Request) {
  sq := &simplanQuery{h}
  h.stdquery(w, r, sq)
}
