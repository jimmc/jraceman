package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type rolepermissionQuery struct{
  h *handler
}

func (sc *rolepermissionQuery) EntityTypeName() string {
  return "rolepermission"
}

func (sc *rolepermissionQuery) NewEntity() interface{} {
  return &domain.RolePermission{}
}

func (sc *rolepermissionQuery) SummaryQuery(format string) string {
  return "select rolepermission.ID as ID, "+
          "role.Name || '[' || role.ID || '] has ' || " +
          " permission.Name || '[' || permission.ID || ']' as summary " +
          "from rolepermission join role on rolepermission.roleid = role.id" +
          " join permission on rolepermission.permissionid = permission.id"
}

func (h *handler) rolepermission(w http.ResponseWriter, r *http.Request) {
  sq := &rolepermissionQuery{h}
  h.stdquery(w, r, sq)
}
