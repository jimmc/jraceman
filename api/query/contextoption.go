package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type contextoptionQuery struct{
  domain.ContextOptionMeta
  h *handler
}

func (sc *contextoptionQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) contextoption(w http.ResponseWriter, r *http.Request) {
  sq := &contextoptionQuery{domain.ContextOptionMeta{}, h}
  h.stdquery(w, r, sq)
}
