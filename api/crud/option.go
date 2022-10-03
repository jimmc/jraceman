package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type optionCrud struct{
  h *handler
}

func (sc *optionCrud) EntityTypeName() string {
  return "option"
}

func (sc *optionCrud) NewEntity() interface{} {
  return &domain.Option{}
}

func (sc *optionCrud) Save(entity interface{}) (string, error) {
  var option *domain.Option = entity.(*domain.Option)
  return sc.h.config.DomainRepos.Option().Save(option)
}

func (sc *optionCrud) List(offset, limit int) ([]interface{}, error) {
  options, err := sc.h.config.DomainRepos.Option().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(options))
  for i, option := range options {
    a[i] = option
  }
  return a, nil
}

func (sc *optionCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Option().FindByID(ID)
}

func (sc *optionCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Option().DeleteByID(ID)
}

func (sc *optionCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldOption *domain.Option = oldEntity.(*domain.Option)
  var newOption *domain.Option = newEntity.(*domain.Option)
  return sc.h.config.DomainRepos.Option().UpdateByID(ID, oldOption, newOption, diffs)
}

func (h *handler) option(w http.ResponseWriter, r *http.Request) {
  sc := &optionCrud{h}
  h.stdcrud(w, r, sc)
}
