package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type exceptionCrud struct{
  h *handler
}

func (sc *exceptionCrud) EntityTypeName() string {
  return "exception"
}

func (sc *exceptionCrud) NewEntity() interface{} {
  return &domain.Exception{}
}

func (sc *exceptionCrud) Save(entity interface{}) (string, error) {
  var exception *domain.Exception = entity.(*domain.Exception)
  return sc.h.config.DomainRepos.Exception().Save(exception)
}

func (sc *exceptionCrud) List(offset, limit int) ([]interface{}, error) {
  exceptions, err := sc.h.config.DomainRepos.Exception().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(exceptions))
  for i, exception := range exceptions {
    a[i] = exception
  }
  return a, nil
}

func (sc *exceptionCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Exception().FindByID(ID)
}

func (sc *exceptionCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Exception().DeleteByID(ID)
}

func (sc *exceptionCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldException *domain.Exception = oldEntity.(*domain.Exception)
  var newException *domain.Exception = newEntity.(*domain.Exception)
  return sc.h.config.DomainRepos.Exception().UpdateByID(ID, oldException, newException, diffs)
}

func (h *handler) exception(w http.ResponseWriter, r *http.Request) {
  sc := &exceptionCrud{h}
  h.stdcrud(w, r, sc)
}
