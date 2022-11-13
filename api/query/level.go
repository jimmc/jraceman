package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type levelQuery struct{
  domain.LevelMeta
  h *handler
}

func (sc *levelQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) level(w http.ResponseWriter, r *http.Request) {
  sq := &levelQuery{domain.LevelMeta{}, h}
  h.stdquery(w, r, sq)
}
