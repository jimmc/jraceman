package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type meetQuery struct{
  domain.MeetMeta
  h *handler
}

func (sc *meetQuery) SummaryQuery(format string) string {
  return "select ID, ShortName || ': ' || Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) meet(w http.ResponseWriter, r *http.Request) {
  sq := &meetQuery{domain.MeetMeta{}, h}
  h.stdquery(w, r, sq)
}
