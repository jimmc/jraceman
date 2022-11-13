package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type raceCrud struct{
  domain.RaceMeta
  h *handler
}

func (sc *raceCrud) Save(entity interface{}) (string, error) {
  var race *domain.Race = entity.(*domain.Race)
  return sc.h.config.DomainRepos.Race().Save(race)
}

func (sc *raceCrud) List(offset, limit int) ([]interface{}, error) {
  races, err := sc.h.config.DomainRepos.Race().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(races))
  for i, race := range races {
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
  sc := &raceCrud{domain.RaceMeta{}, h}
  h.stdcrud(w, r, sc)
}
