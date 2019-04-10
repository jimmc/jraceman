package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type competitionQuery struct{
  h *handler
}

func (sc *competitionQuery) EntityTypeName() string {
  return "competition"
}

func (sc *competitionQuery) NewEntity() interface{} {
  return &domain.Competition{}
}

func (sc *competitionQuery) SummaryQuery() string {
  return "select Name || '[' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) competition(w http.ResponseWriter, r *http.Request) {
  sq := &competitionQuery{h}
  h.stdquery(w, r, sq)
}
