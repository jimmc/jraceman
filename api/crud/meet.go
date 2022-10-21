package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type meetCrud struct{
  h *handler
}

func (sc *meetCrud) EntityTypeName() string {
  return "meet"
}

func (sc *meetCrud) NewEntity() interface{} {
  return &domain.Meet{}
}

func (sc *meetCrud) Save(entity interface{}) (string, error) {
  var meet *domain.Meet = entity.(*domain.Meet)
  return sc.h.config.DomainRepos.Meet().Save(meet)
}

func (sc *meetCrud) List(offset, limit int) ([]interface{}, error) {
  meets, err := sc.h.config.DomainRepos.Meet().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(meets))
  for i, meet := range meets {
    a[i] = meet
  }
  return a, nil
}

func (sc *meetCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Meet().FindByID(ID)
}

func (sc *meetCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Meet().DeleteByID(ID)
}

func (sc *meetCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldMeet *domain.Meet = oldEntity.(*domain.Meet)
  var newMeet *domain.Meet = newEntity.(*domain.Meet)
  return sc.h.config.DomainRepos.Meet().UpdateByID(ID, oldMeet, newMeet, diffs)
}

func (h *handler) meet(w http.ResponseWriter, r *http.Request) {
  sc := &meetCrud{h}
  h.stdcrud(w, r, sc)
}
