package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type seedingplanQuery struct{
  domain.SeedingPlanMeta
  h *handler
}

func (sc *seedingplanQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) seedingplan(w http.ResponseWriter, r *http.Request) {
  sq := &seedingplanQuery{domain.SeedingPlanMeta{}, h}
  h.stdquery(w, r, sq)
}
