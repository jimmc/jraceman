package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type userCrud struct{
  domain.UserMeta
  h *handler
}

func (sc *userCrud) Save(entity interface{}) (string, error) {
  var user *domain.User = entity.(*domain.User)
  return sc.h.config.DomainRepos.User().Save(user)
}

func (sc *userCrud) List(offset, limit int) ([]interface{}, error) {
  users, err := sc.h.config.DomainRepos.User().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(users))
  for i, user := range users {
    a[i] = user
  }
  return a, nil
}

func (sc *userCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.User().FindByID(ID)
}

func (sc *userCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.User().DeleteByID(ID)
}

func (sc *userCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldUser *domain.User = oldEntity.(*domain.User)
  var newUser *domain.User = newEntity.(*domain.User)
  return sc.h.config.DomainRepos.User().UpdateByID(ID, oldUser, newUser, diffs)
}

func (h *handler) user(w http.ResponseWriter, r *http.Request) {
  sc := &userCrud{domain.UserMeta{}, h}
  h.stdcrud(w, r, sc)
}
