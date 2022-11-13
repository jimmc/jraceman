package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type permissionQuery struct{
  domain.PermissionMeta
  h *handler
}

func (sc *permissionQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) permission(w http.ResponseWriter, r *http.Request) {
  sq := &permissionQuery{domain.PermissionMeta{}, h}
  h.stdquery(w, r, sq)
}
