package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type scoringruleCrud struct{
  h *handler
}

func (sc *scoringruleCrud) EntityTypeName() string {
  return "scoringrule"
}

func (sc *scoringruleCrud) NewEntity() interface{} {
  return &domain.ScoringRule{}
}

func (sc *scoringruleCrud) Save(entity interface{}) (string, error) {
  var scoringrule *domain.ScoringRule = entity.(*domain.ScoringRule)
  return sc.h.config.DomainRepos.ScoringRule().Save(scoringrule)
}

func (sc *scoringruleCrud) List(offset, limit int) ([]interface{}, error) {
  scoringrules, err := sc.h.config.DomainRepos.ScoringRule().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(scoringrules))
  for i, scoringrule := range scoringrules {
    a[i] = scoringrule
  }
  return a, nil
}

func (sc *scoringruleCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.ScoringRule().FindByID(ID)
}

func (sc *scoringruleCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.ScoringRule().DeleteByID(ID)
}

func (sc *scoringruleCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldScoringRule *domain.ScoringRule = oldEntity.(*domain.ScoringRule)
  var newScoringRule *domain.ScoringRule = newEntity.(*domain.ScoringRule)
  return sc.h.config.DomainRepos.ScoringRule().UpdateByID(ID, oldScoringRule, newScoringRule, diffs)
}

func (h *handler) scoringrule(w http.ResponseWriter, r *http.Request) {
  sc := &scoringruleCrud{h}
  h.stdcrud(w, r, sc)
}
