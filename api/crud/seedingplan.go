package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type seedingplanCrud struct{
  h *handler
}

func (sc *seedingplanCrud) EntityTypeName() string {
  return "seedingplan"
}

func (sc *seedingplanCrud) NewEntity() interface{} {
  return &domain.SeedingPlan{}
}

func (sc *seedingplanCrud) Save(entity interface{}) (string, error) {
  var seedingplan *domain.SeedingPlan = entity.(*domain.SeedingPlan)
  return sc.h.config.DomainRepos.SeedingPlan().Save(seedingplan)
}

func (sc *seedingplanCrud) List(offset, limit int) ([]interface{}, error) {
  seedingplans, err := sc.h.config.DomainRepos.SeedingPlan().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(seedingplans))
  for i, seedingplan := range seedingplans {
    a[i] = seedingplan
  }
  return a, nil
}

func (sc *seedingplanCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.SeedingPlan().FindByID(ID)
}

func (sc *seedingplanCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.SeedingPlan().DeleteByID(ID)
}

func (sc *seedingplanCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldSeedingPlan *domain.SeedingPlan = oldEntity.(*domain.SeedingPlan)
  var newSeedingPlan *domain.SeedingPlan = newEntity.(*domain.SeedingPlan)
  return sc.h.config.DomainRepos.SeedingPlan().UpdateByID(ID, oldSeedingPlan, newSeedingPlan, diffs)
}

func (h *handler) seedingplan(w http.ResponseWriter, r *http.Request) {
  sc := &seedingplanCrud{h}
  h.stdcrud(w, r, sc)
}
