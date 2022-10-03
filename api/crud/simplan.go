package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type simplanCrud struct{
  h *handler
}

func (sc *simplanCrud) EntityTypeName() string {
  return "simplan"
}

func (sc *simplanCrud) NewEntity() interface{} {
  return &domain.Simplan{}
}

func (sc *simplanCrud) Save(entity interface{}) (string, error) {
  var simplan *domain.Simplan = entity.(*domain.Simplan)
  return sc.h.config.DomainRepos.Simplan().Save(simplan)
}

func (sc *simplanCrud) List(offset, limit int) ([]interface{}, error) {
  simplans, err := sc.h.config.DomainRepos.Simplan().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(simplans))
  for i, simplan := range simplans {
    a[i] = simplan
  }
  return a, nil
}

func (sc *simplanCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Simplan().FindByID(ID)
}

func (sc *simplanCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Simplan().DeleteByID(ID)
}

func (sc *simplanCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldSimplan *domain.Simplan = oldEntity.(*domain.Simplan)
  var newSimplan *domain.Simplan = newEntity.(*domain.Simplan)
  return sc.h.config.DomainRepos.Simplan().UpdateByID(ID, oldSimplan, newSimplan, diffs)
}

func (h *handler) simplan(w http.ResponseWriter, r *http.Request) {
  sc := &simplanCrud{h}
  h.stdcrud(w, r, sc)
}
