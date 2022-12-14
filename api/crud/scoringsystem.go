package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type scoringsystemCrud struct{
  domain.ScoringSystemMeta
  h *handler
}

func (sc *scoringsystemCrud) Save(entity interface{}) (string, error) {
  var scoringsystem *domain.ScoringSystem = entity.(*domain.ScoringSystem)
  return sc.h.config.DomainRepos.ScoringSystem().Save(scoringsystem)
}

func (sc *scoringsystemCrud) List(offset, limit int) ([]interface{}, error) {
  scoringsystems, err := sc.h.config.DomainRepos.ScoringSystem().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(scoringsystems))
  for i, scoringsystem := range scoringsystems {
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
  sc := &scoringsystemCrud{domain.ScoringSystemMeta{}, h}
  h.stdcrud(w, r, sc)
}
