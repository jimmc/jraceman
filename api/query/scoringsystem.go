package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type scoringsystemQuery struct{
  domain.ScoringSystemMeta
  h *handler
}

func (sc *scoringsystemQuery) SummaryQuery(format string) string {
  return "select ID, Name || '[' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) scoringsystem(w http.ResponseWriter, r *http.Request) {
  sq := &scoringsystemQuery{domain.ScoringSystemMeta{}, h}
  h.stdquery(w, r, sq)
}
