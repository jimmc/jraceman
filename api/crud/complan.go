package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type complanCrud struct{
  h *handler
}

func (sc *complanCrud) EntityTypeName() string {
  return "complan"
}

func (sc *complanCrud) NewEntity() interface{} {
  return &domain.Complan{}
}

func (sc *complanCrud) Save(entity interface{}) (string, error) {
  var complan *domain.Complan = entity.(*domain.Complan)
  return sc.h.config.DomainRepos.Complan().Save(complan)
}

func (sc *complanCrud) List(offset, limit int) ([]interface{}, error) {
  sites, err := sc.h.config.DomainRepos.Complan().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(sites))
  for i, complan := range sites {
    a[i] = complan
  }
  return a, nil
}

func (sc *complanCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Complan().FindByID(ID)
}

func (sc *complanCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Complan().DeleteByID(ID)
}

func (sc *complanCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldComplan *domain.Complan = oldEntity.(*domain.Complan)
  var newComplan *domain.Complan = newEntity.(*domain.Complan)
  return sc.h.config.DomainRepos.Complan().UpdateByID(ID, oldComplan, newComplan, diffs)
}

func (h *handler) complan(w http.ResponseWriter, r *http.Request) {
  sc := &complanCrud{h}
  h.stdcrud(w, r, sc)
}
