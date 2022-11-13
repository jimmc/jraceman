package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type roleroleCrud struct{
  domain.RoleRoleMeta
  h *handler
}

func (sc *roleroleCrud) Save(entity interface{}) (string, error) {
  var rolerole *domain.RoleRole = entity.(*domain.RoleRole)
  return sc.h.config.DomainRepos.RoleRole().Save(rolerole)
}

func (sc *roleroleCrud) List(offset, limit int) ([]interface{}, error) {
  roleroles, err := sc.h.config.DomainRepos.RoleRole().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(roleroles))
  for i, rolerole := range roleroles {
    a[i] = rolerole
  }
  return a, nil
}

func (sc *roleroleCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.RoleRole().FindByID(ID)
}

func (sc *roleroleCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.RoleRole().DeleteByID(ID)
}

func (sc *roleroleCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldRoleRole *domain.RoleRole = oldEntity.(*domain.RoleRole)
  var newRoleRole *domain.RoleRole = newEntity.(*domain.RoleRole)
  return sc.h.config.DomainRepos.RoleRole().UpdateByID(ID, oldRoleRole, newRoleRole, diffs)
}

func (h *handler) rolerole(w http.ResponseWriter, r *http.Request) {
  sc := &roleroleCrud{domain.RoleRoleMeta{}, h}
  h.stdcrud(w, r, sc)
}
