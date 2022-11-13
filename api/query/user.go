package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type userQuery struct{
  domain.UserMeta
  h *handler
}

func (sc *userQuery) SummaryQuery(format string) string {
  return "select ID, Username || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) user(w http.ResponseWriter, r *http.Request) {
  sq := &userQuery{domain.UserMeta{}, h}
  h.stdquery(w, r, sq)
}
