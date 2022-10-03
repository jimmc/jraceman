package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type simplanruleCrud struct{
  h *handler
}

func (sc *simplanruleCrud) EntityTypeName() string {
  return "simplanrule"
}

func (sc *simplanruleCrud) NewEntity() interface{} {
  return &domain.SimplanRule{}
}

func (sc *simplanruleCrud) Save(entity interface{}) (string, error) {
  var simplanrule *domain.SimplanRule = entity.(*domain.SimplanRule)
  return sc.h.config.DomainRepos.SimplanRule().Save(simplanrule)
}

func (sc *simplanruleCrud) List(offset, limit int) ([]interface{}, error) {
  simplanrules, err := sc.h.config.DomainRepos.SimplanRule().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(simplanrules))
  for i, simplanrule := range simplanrules {
    a[i] = simplanrule
  }
  return a, nil
}

func (sc *simplanruleCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.SimplanRule().FindByID(ID)
}

func (sc *simplanruleCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.SimplanRule().DeleteByID(ID)
}

func (sc *simplanruleCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldSimplanRule *domain.SimplanRule = oldEntity.(*domain.SimplanRule)
  var newSimplanRule *domain.SimplanRule = newEntity.(*domain.SimplanRule)
  return sc.h.config.DomainRepos.SimplanRule().UpdateByID(ID, oldSimplanRule, newSimplanRule, diffs)
}

func (h *handler) simplanrule(w http.ResponseWriter, r *http.Request) {
  sc := &simplanruleCrud{h}
  h.stdcrud(w, r, sc)
}
