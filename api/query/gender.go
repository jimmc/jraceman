package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type genderQuery struct{
  domain.GenderMeta
  h *handler
}

func (sc *genderQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) gender(w http.ResponseWriter, r *http.Request) {
  sq := &genderQuery{domain.GenderMeta{}, h}
  h.stdquery(w, r, sq)
}
