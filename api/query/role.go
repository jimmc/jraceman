package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type roleQuery struct{
  h *handler
}

func (sc *roleQuery) EntityTypeName() string {
  return "role"
}

func (sc *roleQuery) NewEntity() interface{} {
  return &domain.Role{}
}

func (sc *roleQuery) SummaryQuery() string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) role(w http.ResponseWriter, r *http.Request) {
  sq := &roleQuery{h}
  h.stdquery(w, r, sq)
}
