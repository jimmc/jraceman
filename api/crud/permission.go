package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type permissionCrud struct{
  h *handler
}

func (sc *permissionCrud) EntityTypeName() string {
  return "permission"
}

func (sc *permissionCrud) NewEntity() interface{} {
  return &domain.Permission{}
}

func (sc *permissionCrud) Save(entity interface{}) (string, error) {
  var permission *domain.Permission = entity.(*domain.Permission)
  return sc.h.config.DomainRepos.Permission().Save(permission)
}

func (sc *permissionCrud) List(offset, limit int) ([]interface{}, error) {
  permissions, err := sc.h.config.DomainRepos.Permission().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(permissions))
  for i, permission := range permissions {
    a[i] = permission
  }
  return a, nil
}

func (sc *permissionCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Permission().FindByID(ID)
}

func (sc *permissionCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Permission().DeleteByID(ID)
}

func (sc *permissionCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldPermission *domain.Permission = oldEntity.(*domain.Permission)
  var newPermission *domain.Permission = newEntity.(*domain.Permission)
  return sc.h.config.DomainRepos.Permission().UpdateByID(ID, oldPermission, newPermission, diffs)
}

func (h *handler) permission(w http.ResponseWriter, r *http.Request) {
  sc := &permissionCrud{h}
  h.stdcrud(w, r, sc)
}
