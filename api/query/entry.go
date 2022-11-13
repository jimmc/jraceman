package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type entryQuery struct{
  domain.EntryMeta
  h *handler
}

func (sc *entryQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) entry(w http.ResponseWriter, r *http.Request) {
  sq := &entryQuery{domain.EntryMeta{}, h}
  h.stdquery(w, r, sq)
}
