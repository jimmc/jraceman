package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type stageQuery struct{
  domain.StageMeta
  h *handler
}

func (sc *stageQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) stage(w http.ResponseWriter, r *http.Request) {
  sq := &stageQuery{domain.StageMeta{}, h}
  h.stdquery(w, r, sq)
}
