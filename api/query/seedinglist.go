package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type seedinglistQuery struct{
  domain.SeedingListMeta
  h *handler
}

func (sc *seedinglistQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) seedinglist(w http.ResponseWriter, r *http.Request) {
  sq := &seedinglistQuery{domain.SeedingListMeta{}, h}
  h.stdquery(w, r, sq)
}
