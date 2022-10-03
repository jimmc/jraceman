package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type roleCrud struct{
  h *handler
}

func (sc *roleCrud) EntityTypeName() string {
  return "role"
}

func (sc *roleCrud) NewEntity() interface{} {
  return &domain.Role{}
}

func (sc *roleCrud) Save(entity interface{}) (string, error) {
  var role *domain.Role = entity.(*domain.Role)
  return sc.h.config.DomainRepos.Role().Save(role)
}

func (sc *roleCrud) List(offset, limit int) ([]interface{}, error) {
  roles, err := sc.h.config.DomainRepos.Role().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(roles))
  for i, role := range roles {
    a[i] = role
  }
  return a, nil
}

func (sc *roleCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Role().FindByID(ID)
}

func (sc *roleCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Role().DeleteByID(ID)
}

func (sc *roleCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldRole *domain.Role = oldEntity.(*domain.Role)
  var newRole *domain.Role = newEntity.(*domain.Role)
  return sc.h.config.DomainRepos.Role().UpdateByID(ID, oldRole, newRole, diffs)
}

func (h *handler) role(w http.ResponseWriter, r *http.Request) {
  sc := &roleCrud{h}
  h.stdcrud(w, r, sc)
}
