package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type meetQuery struct{
  h *handler
}

func (sc *meetQuery) EntityTypeName() string {
  return "meet"
}

func (sc *meetQuery) NewEntity() interface{} {
  return &domain.Meet{}
}

func (sc *meetQuery) SummaryQuery(format string) string {
  return "select ID, ShortName || ': ' || Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) meet(w http.ResponseWriter, r *http.Request) {
  sq := &meetQuery{h}
  h.stdquery(w, r, sq)
}
