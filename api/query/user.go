package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type userQuery struct{
  h *handler
}

func (sc *userQuery) EntityTypeName() string {
  return "user"
}

func (sc *userQuery) NewEntity() interface{} {
  return &domain.User{}
}

func (sc *userQuery) SummaryQuery(format string) string {
  return "select ID, Username || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) user(w http.ResponseWriter, r *http.Request) {
  sq := &userQuery{h}
  h.stdquery(w, r, sq)
}
