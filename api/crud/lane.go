package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type laneCrud struct{
  h *handler
}

func (sc *laneCrud) EntityTypeName() string {
  return "lane"
}

func (sc *laneCrud) NewEntity() interface{} {
  return &domain.Lane{}
}

func (sc *laneCrud) Save(entity interface{}) (string, error) {
  var lane *domain.Lane = entity.(*domain.Lane)
  return sc.h.config.DomainRepos.Lane().Save(lane)
}

func (sc *laneCrud) List(offset, limit int) ([]interface{}, error) {
  lanes, err := sc.h.config.DomainRepos.Lane().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(lanes))
  for i, lane := range lanes {
    a[i] = lane
  }
  return a, nil
}

func (sc *laneCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Lane().FindByID(ID)
}

func (sc *laneCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Lane().DeleteByID(ID)
}

func (sc *laneCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldLane *domain.Lane = oldEntity.(*domain.Lane)
  var newLane *domain.Lane = newEntity.(*domain.Lane)
  return sc.h.config.DomainRepos.Lane().UpdateByID(ID, oldLane, newLane, diffs)
}

func (h *handler) lane(w http.ResponseWriter, r *http.Request) {
  sc := &laneCrud{h}
  h.stdcrud(w, r, sc)
}
