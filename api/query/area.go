package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type areaQuery struct{
  h *handler
}

func (sc *areaQuery) EntityTypeName() string {
  return "area"
}

func (sc *areaQuery) NewEntity() interface{} {
  return &domain.Area{}
}

func (sc *areaQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) area(w http.ResponseWriter, r *http.Request) {
  sq := &areaQuery{h}
  h.stdquery(w, r, sq)
}
