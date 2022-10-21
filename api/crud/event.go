package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type eventCrud struct{
  h *handler
}

func (sc *eventCrud) EntityTypeName() string {
  return "event"
}

func (sc *eventCrud) NewEntity() interface{} {
  return &domain.Event{}
}

func (sc *eventCrud) Save(entity interface{}) (string, error) {
  var event *domain.Event = entity.(*domain.Event)
  return sc.h.config.DomainRepos.Event().Save(event)
}

func (sc *eventCrud) List(offset, limit int) ([]interface{}, error) {
  events, err := sc.h.config.DomainRepos.Event().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(events))
  for i, event := range events {
    a[i] = event
  }
  return a, nil
}

func (sc *eventCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Event().FindByID(ID)
}

func (sc *eventCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Event().DeleteByID(ID)
}

func (sc *eventCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldEvent *domain.Event = oldEntity.(*domain.Event)
  var newEvent *domain.Event = newEntity.(*domain.Event)
  return sc.h.config.DomainRepos.Event().UpdateByID(ID, oldEvent, newEvent, diffs)
}

func (h *handler) event(w http.ResponseWriter, r *http.Request) {
  sc := &eventCrud{h}
  h.stdcrud(w, r, sc)
}
