package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type progressionQuery struct{
  domain.ProgressionMeta
  h *handler
}

func (sc *progressionQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) progression(w http.ResponseWriter, r *http.Request) {
  sq := &progressionQuery{domain.ProgressionMeta{}, h}
  h.stdquery(w, r, sq)
}
