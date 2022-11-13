package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type roleroleQuery struct{
  domain.RoleRoleMeta
  h *handler
}

func (sc *roleroleQuery) SummaryQuery(format string) string {
  return "select rolerole.ID as ID, "+
          "role.Name || '[' || role.ID || '] has ' || " +
          " hasrole.Name || '[' || hasrole.ID || ']' as summary " +
          "from rolerole join role on rolerole.roleid = role.id" +
          " join role as hasrole on rolerole.hasroleid = hasrole.id"
}

func (h *handler) rolerole(w http.ResponseWriter, r *http.Request) {
  sq := &roleroleQuery{domain.RoleRoleMeta{}, h}
  h.stdquery(w, r, sq)
}
