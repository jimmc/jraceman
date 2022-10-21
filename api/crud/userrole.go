package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type userroleCrud struct{
  h *handler
}

func (sc *userroleCrud) EntityTypeName() string {
  return "userrole"
}

func (sc *userroleCrud) NewEntity() interface{} {
  return &domain.UserRole{}
}

func (sc *userroleCrud) Save(entity interface{}) (string, error) {
  var userrole *domain.UserRole = entity.(*domain.UserRole)
  return sc.h.config.DomainRepos.UserRole().Save(userrole)
}

func (sc *userroleCrud) List(offset, limit int) ([]interface{}, error) {
  userroles, err := sc.h.config.DomainRepos.UserRole().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(userroles))
  for i, userrole := range userroles {
    a[i] = userrole
  }
  return a, nil
}

func (sc *userroleCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.UserRole().FindByID(ID)
}

func (sc *userroleCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.UserRole().DeleteByID(ID)
}

func (sc *userroleCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldUserRole *domain.UserRole = oldEntity.(*domain.UserRole)
  var newUserRole *domain.UserRole = newEntity.(*domain.UserRole)
  return sc.h.config.DomainRepos.UserRole().UpdateByID(ID, oldUserRole, newUserRole, diffs)
}

func (h *handler) userrole(w http.ResponseWriter, r *http.Request) {
  sc := &userroleCrud{h}
  h.stdcrud(w, r, sc)
}
