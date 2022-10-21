package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type seedinglistCrud struct{
  h *handler
}

func (sc *seedinglistCrud) EntityTypeName() string {
  return "seedinglist"
}

func (sc *seedinglistCrud) NewEntity() interface{} {
  return &domain.SeedingList{}
}

func (sc *seedinglistCrud) Save(entity interface{}) (string, error) {
  var seedinglist *domain.SeedingList = entity.(*domain.SeedingList)
  return sc.h.config.DomainRepos.SeedingList().Save(seedinglist)
}

func (sc *seedinglistCrud) List(offset, limit int) ([]interface{}, error) {
  seedinglists, err := sc.h.config.DomainRepos.SeedingList().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(seedinglists))
  for i, seedinglist := range seedinglists {
    a[i] = seedinglist
  }
  return a, nil
}

func (sc *seedinglistCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.SeedingList().FindByID(ID)
}

func (sc *seedinglistCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.SeedingList().DeleteByID(ID)
}

func (sc *seedinglistCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldSeedingList *domain.SeedingList = oldEntity.(*domain.SeedingList)
  var newSeedingList *domain.SeedingList = newEntity.(*domain.SeedingList)
  return sc.h.config.DomainRepos.SeedingList().UpdateByID(ID, oldSeedingList, newSeedingList, diffs)
}

func (h *handler) seedinglist(w http.ResponseWriter, r *http.Request) {
  sc := &seedinglistCrud{h}
  h.stdcrud(w, r, sc)
}
