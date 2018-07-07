package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type raceCrud struct{
  h *handler
}

func (sc *raceCrud) EntityTypeName() string {
  return "race"
}

func (sc *raceCrud) NewEntity() interface{} {
  return &domain.Race{}
}

func (sc *raceCrud) Save(entity interface{}) (string, error) {
  var race *domain.Race = entity.(*domain.Race)
  return sc.h.config.DomainRepos.Race().Save(race)
}

func (sc *raceCrud) List(offset, limit int) ([]interface{}, error) {
  sites, err := sc.h.config.DomainRepos.Race().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(sites))
  for i, race := range sites {
    a[i] = race
  }
  return a, nil
}

func (sc *raceCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Race().FindByID(ID)
}

func (sc *raceCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Race().DeleteByID(ID)
}

func (sc *raceCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldRace *domain.Race = oldEntity.(*domain.Race)
  var newRace *domain.Race = newEntity.(*domain.Race)
  return sc.h.config.DomainRepos.Race().UpdateByID(ID, oldRace, newRace, diffs)
}

func (h *handler) race(w http.ResponseWriter, r *http.Request) {
  sc := &raceCrud{h}
  h.stdcrud(w, r, sc)
}
