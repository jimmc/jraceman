package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type permissionQuery struct{
  h *handler
}

func (sc *permissionQuery) EntityTypeName() string {
  return "permission"
}

func (sc *permissionQuery) NewEntity() interface{} {
  return &domain.Permission{}
}

func (sc *permissionQuery) SummaryQuery() string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) permission(w http.ResponseWriter, r *http.Request) {
  sq := &permissionQuery{h}
  h.stdquery(w, r, sq)
}
