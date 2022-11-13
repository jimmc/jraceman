package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type simplanstageQuery struct{
  domain.SimplanStageMeta
  h *handler
}

func (sc *simplanstageQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) simplanstage(w http.ResponseWriter, r *http.Request) {
  sq := &simplanstageQuery{domain.SimplanStageMeta{}, h}
  h.stdquery(w, r, sq)
}
