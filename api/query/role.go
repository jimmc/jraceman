package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type roleQuery struct{
  domain.RoleMeta
  h *handler
}

func (sc *roleQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) role(w http.ResponseWriter, r *http.Request) {
  sq := &roleQuery{domain.RoleMeta{}, h}
  h.stdquery(w, r, sq)
}
