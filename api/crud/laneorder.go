package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type laneorderCrud struct{
  domain.LaneOrderMeta
  h *handler
}

func (sc *laneorderCrud) Save(entity interface{}) (string, error) {
  var laneorder *domain.LaneOrder = entity.(*domain.LaneOrder)
  return sc.h.config.DomainRepos.LaneOrder().Save(laneorder)
}

func (sc *laneorderCrud) List(offset, limit int) ([]interface{}, error) {
  laneorders, err := sc.h.config.DomainRepos.LaneOrder().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(laneorders))
  for i, laneorder := range laneorders {
    a[i] = laneorder
  }
  return a, nil
}

func (sc *laneorderCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.LaneOrder().FindByID(ID)
}

func (sc *laneorderCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.LaneOrder().DeleteByID(ID)
}

func (sc *laneorderCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldLaneOrder *domain.LaneOrder = oldEntity.(*domain.LaneOrder)
  var newLaneOrder *domain.LaneOrder = newEntity.(*domain.LaneOrder)
  return sc.h.config.DomainRepos.LaneOrder().UpdateByID(ID, oldLaneOrder, newLaneOrder, diffs)
}

func (h *handler) laneorder(w http.ResponseWriter, r *http.Request) {
  sc := &laneorderCrud{domain.LaneOrderMeta{}, h}
  h.stdcrud(w, r, sc)
}
