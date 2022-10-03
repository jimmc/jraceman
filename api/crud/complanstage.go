package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type complanstageCrud struct{
  h *handler
}

func (sc *complanstageCrud) EntityTypeName() string {
  return "complanstage"
}

func (sc *complanstageCrud) NewEntity() interface{} {
  return &domain.ComplanStage{}
}

func (sc *complanstageCrud) Save(entity interface{}) (string, error) {
  var complanstage *domain.ComplanStage = entity.(*domain.ComplanStage)
  return sc.h.config.DomainRepos.ComplanStage().Save(complanstage)
}

func (sc *complanstageCrud) List(offset, limit int) ([]interface{}, error) {
  complanstages, err := sc.h.config.DomainRepos.ComplanStage().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(complanstages))
  for i, complanstage := range complanstages {
    a[i] = complanstage
  }
  return a, nil
}

func (sc *complanstageCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.ComplanStage().FindByID(ID)
}

func (sc *complanstageCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.ComplanStage().DeleteByID(ID)
}

func (sc *complanstageCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldComplanStage *domain.ComplanStage = oldEntity.(*domain.ComplanStage)
  var newComplanStage *domain.ComplanStage = newEntity.(*domain.ComplanStage)
  return sc.h.config.DomainRepos.ComplanStage().UpdateByID(ID, oldComplanStage, newComplanStage, diffs)
}

func (h *handler) complanstage(w http.ResponseWriter, r *http.Request) {
  sc := &complanstageCrud{h}
  h.stdcrud(w, r, sc)
}
