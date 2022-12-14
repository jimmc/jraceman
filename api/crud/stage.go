package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type stageCrud struct{
  domain.StageMeta
  h *handler
}

func (sc *stageCrud) Save(entity interface{}) (string, error) {
  var stage *domain.Stage = entity.(*domain.Stage)
  return sc.h.config.DomainRepos.Stage().Save(stage)
}

func (sc *stageCrud) List(offset, limit int) ([]interface{}, error) {
  stages, err := sc.h.config.DomainRepos.Stage().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(stages))
  for i, stage := range stages {
    a[i] = stage
  }
  return a, nil
}

func (sc *stageCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Stage().FindByID(ID)
}

func (sc *stageCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Stage().DeleteByID(ID)
}

func (sc *stageCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldStage *domain.Stage = oldEntity.(*domain.Stage)
  var newStage *domain.Stage = newEntity.(*domain.Stage)
  return sc.h.config.DomainRepos.Stage().UpdateByID(ID, oldStage, newStage, diffs)
}

func (h *handler) stage(w http.ResponseWriter, r *http.Request) {
  sc := &stageCrud{domain.StageMeta{}, h}
  h.stdcrud(w, r, sc)
}
