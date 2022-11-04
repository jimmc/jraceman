package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
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

func (sc *competitionQuery) SummaryQuery(format string) string {
  return "select ID, Name || '[' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) competition(w http.ResponseWriter, r *http.Request) {
  sq := &competitionQuery{h}
  h.stdquery(w, r, sq)
}
