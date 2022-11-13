package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type areaQuery struct{
  domain.AreaMeta
  h *handler
}

func (sc *areaQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) area(w http.ResponseWriter, r *http.Request) {
  sq := &areaQuery{domain.AreaMeta{}, h}
  h.stdquery(w, r, sq)
}
