package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type userroleQuery struct{
  h *handler
}

func (sc *userroleQuery) EntityTypeName() string {
  return "userrole"
}

func (sc *userroleQuery) NewEntity() interface{} {
  return &domain.UserRole{}
}

func (sc *userroleQuery) SummaryQuery(format string) string {
  return "select userrole.ID as ID, "+
          "user.Username || '[' || user.ID || '] has ' || " +
          " role.Name || '[' || role.ID || ']' as summary " +
          "from userrole join user on userrole.userid = user.id" +
          " join role on userrole.roleid = role.id"
}

func (h *handler) userrole(w http.ResponseWriter, r *http.Request) {
  sq := &userroleQuery{h}
  h.stdquery(w, r, sq)
}
