package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type simplanstageCrud struct{
  h *handler
}

func (sc *simplanstageCrud) EntityTypeName() string {
  return "simplanstage"
}

func (sc *simplanstageCrud) NewEntity() interface{} {
  return &domain.SimplanStage{}
}

func (sc *simplanstageCrud) Save(entity interface{}) (string, error) {
  var simplanstage *domain.SimplanStage = entity.(*domain.SimplanStage)
  return sc.h.config.DomainRepos.SimplanStage().Save(simplanstage)
}

func (sc *simplanstageCrud) List(offset, limit int) ([]interface{}, error) {
  simplanstages, err := sc.h.config.DomainRepos.SimplanStage().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(simplanstages))
  for i, simplanstage := range simplanstages {
    a[i] = simplanstage
  }
  return a, nil
}

func (sc *simplanstageCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.SimplanStage().FindByID(ID)
}

func (sc *simplanstageCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.SimplanStage().DeleteByID(ID)
}

func (sc *simplanstageCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldSimplanStage *domain.SimplanStage = oldEntity.(*domain.SimplanStage)
  var newSimplanStage *domain.SimplanStage = newEntity.(*domain.SimplanStage)
  return sc.h.config.DomainRepos.SimplanStage().UpdateByID(ID, oldSimplanStage, newSimplanStage, diffs)
}

func (h *handler) simplanstage(w http.ResponseWriter, r *http.Request) {
  sc := &simplanstageCrud{h}
  h.stdcrud(w, r, sc)
}
