package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type competitionQuery struct{
  domain.CompetitionMeta
  h *handler
}

func (sc *competitionQuery) SummaryQuery(format string) string {
  return "select ID, Name || '[' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) competition(w http.ResponseWriter, r *http.Request) {
  sq := &competitionQuery{domain.CompetitionMeta{}, h}
  h.stdquery(w, r, sq)
}
