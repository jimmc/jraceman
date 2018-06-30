package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type personCrud struct{
  h *handler
}

func (sc *personCrud) EntityTypeName() string {
  return "person"
}

func (sc *personCrud) NewEntity() interface{} {
  return &domain.Person{}
}

func (sc *personCrud) Save(entity interface{}) (string, error) {
  var person *domain.Person = entity.(*domain.Person)
  return sc.h.config.DomainRepos.Person().Save(person)
}

func (sc *personCrud) List(offset, limit int) ([]interface{}, error) {
  sites, err := sc.h.config.DomainRepos.Person().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(sites))
  for i, person := range sites {
    a[i] = person
  }
  return a, nil
}

func (sc *personCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Person().FindByID(ID)
}

func (sc *personCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Person().DeleteByID(ID)
}

func (sc *personCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldPerson *domain.Person = oldEntity.(*domain.Person)
  var newPerson *domain.Person = newEntity.(*domain.Person)
  return sc.h.config.DomainRepos.Person().UpdateByID(ID, oldPerson, newPerson, diffs)
}

func (h *handler) person(w http.ResponseWriter, r *http.Request) {
  sc := &personCrud{h}
  h.stdcrud(w, r, sc)
}
