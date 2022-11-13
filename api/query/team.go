package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type teamQuery struct{
  domain.TeamMeta
  h *handler
}

func (sc *teamQuery) SummaryQuery(format string) string {
  return "select ID, ShortName || ': ' || Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) team(w http.ResponseWriter, r *http.Request) {
  sq := &teamQuery{domain.TeamMeta{}, h}
  h.stdquery(w, r, sq)
}
