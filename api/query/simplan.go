package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type simplanQuery struct{
  domain.SimplanMeta
  h *handler
}

func (sc *simplanQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) simplan(w http.ResponseWriter, r *http.Request) {
  sq := &simplanQuery{domain.SimplanMeta{}, h}
  h.stdquery(w, r, sq)
}
