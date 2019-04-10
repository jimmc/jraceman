package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
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

func (sc *teamQuery) SummaryQuery() string {
  return "select ShortName || ': ' || Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) team(w http.ResponseWriter, r *http.Request) {
  sq := &teamQuery{h}
  h.stdquery(w, r, sq)
}
