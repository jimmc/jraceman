package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type complanruleCrud struct{
  h *handler
}

func (sc *complanruleCrud) EntityTypeName() string {
  return "complanrule"
}

func (sc *complanruleCrud) NewEntity() interface{} {
  return &domain.ComplanRule{}
}

func (sc *complanruleCrud) Save(entity interface{}) (string, error) {
  var complanrule *domain.ComplanRule = entity.(*domain.ComplanRule)
  return sc.h.config.DomainRepos.ComplanRule().Save(complanrule)
}

func (sc *complanruleCrud) List(offset, limit int) ([]interface{}, error) {
  sites, err := sc.h.config.DomainRepos.ComplanRule().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(sites))
  for i, complanrule := range sites {
    a[i] = complanrule
  }
  return a, nil
}

func (sc *complanruleCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.ComplanRule().FindByID(ID)
}

func (sc *complanruleCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.ComplanRule().DeleteByID(ID)
}

func (sc *complanruleCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldComplanRule *domain.ComplanRule = oldEntity.(*domain.ComplanRule)
  var newComplanRule *domain.ComplanRule = newEntity.(*domain.ComplanRule)
  return sc.h.config.DomainRepos.ComplanRule().UpdateByID(ID, oldComplanRule, newComplanRule, diffs)
}

func (h *handler) complanrule(w http.ResponseWriter, r *http.Request) {
  sc := &complanruleCrud{h}
  h.stdcrud(w, r, sc)
}
