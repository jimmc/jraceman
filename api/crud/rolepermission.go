package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type rolepermissionCrud struct{
  h *handler
}

func (sc *rolepermissionCrud) EntityTypeName() string {
  return "rolepermission"
}

func (sc *rolepermissionCrud) NewEntity() interface{} {
  return &domain.RolePermission{}
}

func (sc *rolepermissionCrud) Save(entity interface{}) (string, error) {
  var rolepermission *domain.RolePermission = entity.(*domain.RolePermission)
  return sc.h.config.DomainRepos.RolePermission().Save(rolepermission)
}

func (sc *rolepermissionCrud) List(offset, limit int) ([]interface{}, error) {
  rolepermissions, err := sc.h.config.DomainRepos.RolePermission().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(rolepermissions))
  for i, rolepermission := range rolepermissions {
    a[i] = rolepermission
  }
  return a, nil
}

func (sc *rolepermissionCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.RolePermission().FindByID(ID)
}

func (sc *rolepermissionCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.RolePermission().DeleteByID(ID)
}

func (sc *rolepermissionCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldRolePermission *domain.RolePermission = oldEntity.(*domain.RolePermission)
  var newRolePermission *domain.RolePermission = newEntity.(*domain.RolePermission)
  return sc.h.config.DomainRepos.RolePermission().UpdateByID(ID, oldRolePermission, newRolePermission, diffs)
}

func (h *handler) rolepermission(w http.ResponseWriter, r *http.Request) {
  sc := &rolepermissionCrud{h}
  h.stdcrud(w, r, sc)
}
