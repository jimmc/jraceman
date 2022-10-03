package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type entryCrud struct{
  h *handler
}

func (sc *entryCrud) EntityTypeName() string {
  return "entry"
}

func (sc *entryCrud) NewEntity() interface{} {
  return &domain.Entry{}
}

func (sc *entryCrud) Save(entity interface{}) (string, error) {
  var entry *domain.Entry = entity.(*domain.Entry)
  return sc.h.config.DomainRepos.Entry().Save(entry)
}

func (sc *entryCrud) List(offset, limit int) ([]interface{}, error) {
  entries, err := sc.h.config.DomainRepos.Entry().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(entries))
  for i, entry := range entries {
    a[i] = entry
  }
  return a, nil
}

func (sc *entryCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Entry().FindByID(ID)
}

func (sc *entryCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Entry().DeleteByID(ID)
}

func (sc *entryCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldEntry *domain.Entry = oldEntity.(*domain.Entry)
  var newEntry *domain.Entry = newEntity.(*domain.Entry)
  return sc.h.config.DomainRepos.Entry().UpdateByID(ID, oldEntry, newEntry, diffs)
}

func (h *handler) entry(w http.ResponseWriter, r *http.Request) {
  sc := &entryCrud{h}
  h.stdcrud(w, r, sc)
}
