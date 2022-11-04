package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type teamQuery struct{
  h *handler
}

func (sc *teamQuery) EntityTypeName() string {
  return "team"
}

func (sc *teamQuery) NewEntity() interface{} {
  return &domain.Team{}
}

func (sc *teamQuery) SummaryQuery(format string) string {
  return "select ID, ShortName || ': ' || Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) team(w http.ResponseWriter, r *http.Request) {
  sq := &teamQuery{h}
  h.stdquery(w, r, sq)
}
