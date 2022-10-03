package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type contextoptionCrud struct{
  h *handler
}

func (sc *contextoptionCrud) EntityTypeName() string {
  return "contextoption"
}

func (sc *contextoptionCrud) NewEntity() interface{} {
  return &domain.ContextOption{}
}

func (sc *contextoptionCrud) Save(entity interface{}) (string, error) {
  var contextoption *domain.ContextOption = entity.(*domain.ContextOption)
  return sc.h.config.DomainRepos.ContextOption().Save(contextoption)
}

func (sc *contextoptionCrud) List(offset, limit int) ([]interface{}, error) {
  contextoptions, err := sc.h.config.DomainRepos.ContextOption().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(contextoptions))
  for i, contextoption := range contextoptions {
    a[i] = contextoption
  }
  return a, nil
}

func (sc *contextoptionCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.ContextOption().FindByID(ID)
}

func (sc *contextoptionCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.ContextOption().DeleteByID(ID)
}

func (sc *contextoptionCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldContextOption *domain.ContextOption = oldEntity.(*domain.ContextOption)
  var newContextOption *domain.ContextOption = newEntity.(*domain.ContextOption)
  return sc.h.config.DomainRepos.ContextOption().UpdateByID(ID, oldContextOption, newContextOption, diffs)
}

func (h *handler) contextoption(w http.ResponseWriter, r *http.Request) {
  sc := &contextoptionCrud{h}
  h.stdcrud(w, r, sc)
}
