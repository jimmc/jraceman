package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type progressionCrud struct{
  h *handler
}

func (sc *progressionCrud) EntityTypeName() string {
  return "progression"
}

func (sc *progressionCrud) NewEntity() interface{} {
  return &domain.Progression{}
}

func (sc *progressionCrud) Save(entity interface{}) (string, error) {
  var progression *domain.Progression = entity.(*domain.Progression)
  return sc.h.config.DomainRepos.Progression().Save(progression)
}

func (sc *progressionCrud) List(offset, limit int) ([]interface{}, error) {
  progressions, err := sc.h.config.DomainRepos.Progression().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(progressions))
  for i, progression := range progressions {
    a[i] = progression
  }
  return a, nil
}

func (sc *progressionCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Progression().FindByID(ID)
}

func (sc *progressionCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Progression().DeleteByID(ID)
}

func (sc *progressionCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldProgression *domain.Progression = oldEntity.(*domain.Progression)
  var newProgression *domain.Progression = newEntity.(*domain.Progression)
  return sc.h.config.DomainRepos.Progression().UpdateByID(ID, oldProgression, newProgression, diffs)
}

func (h *handler) progression(w http.ResponseWriter, r *http.Request) {
  sc := &progressionCrud{h}
  h.stdcrud(w, r, sc)
}
