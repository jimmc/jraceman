package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type scoringsystemCrud struct{
  h *handler
}

func (sc *scoringsystemCrud) EntityTypeName() string {
  return "scoringsystem"
}

func (sc *scoringsystemCrud) NewEntity() interface{} {
  return &domain.ScoringSystem{}
}

func (sc *scoringsystemCrud) Save(entity interface{}) (string, error) {
  var scoringsystem *domain.ScoringSystem = entity.(*domain.ScoringSystem)
  return sc.h.config.DomainRepos.ScoringSystem().Save(scoringsystem)
}

func (sc *scoringsystemCrud) List(offset, limit int) ([]interface{}, error) {
  sites, err := sc.h.config.DomainRepos.ScoringSystem().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(sites))
  for i, scoringsystem := range sites {
    a[i] = scoringsystem
  }
  return a, nil
}

func (sc *scoringsystemCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.ScoringSystem().FindByID(ID)
}

func (sc *scoringsystemCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.ScoringSystem().DeleteByID(ID)
}

func (sc *scoringsystemCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldScoringSystem *domain.ScoringSystem = oldEntity.(*domain.ScoringSystem)
  var newScoringSystem *domain.ScoringSystem = newEntity.(*domain.ScoringSystem)
  return sc.h.config.DomainRepos.ScoringSystem().UpdateByID(ID, oldScoringSystem, newScoringSystem, diffs)
}

func (h *handler) scoringsystem(w http.ResponseWriter, r *http.Request) {
  sc := &scoringsystemCrud{h}
  h.stdcrud(w, r, sc)
}
