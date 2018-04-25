package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type siteCrud struct{
  h *handler
}

func (sc *siteCrud) EntityTypeName() string {
  return "site"
}

func (sc *siteCrud) NewEntity() interface{} {
  return &domain.Site{}
}

func (sc *siteCrud) Save(entity interface{}) error {
  var site *domain.Site = entity.(*domain.Site)
  return sc.h.config.DomainRepos.Site().Save(site)
}

func (sc *siteCrud) FindById(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Site().FindById(ID)
}

func (sc *siteCrud) DeleteById(ID string) error {
  return sc.h.config.DomainRepos.Site().DeleteById(ID)
}

func (sc *siteCrud) UpdateById(ID string, oldEntity, newEntity interface{}) error {
  var oldSite *domain.Site = oldEntity.(*domain.Site)
  var newSite *domain.Site = newEntity.(*domain.Site)
  return sc.h.config.DomainRepos.Site().UpdateById(ID, oldSite, newSite)
}

func (h *handler) site(w http.ResponseWriter, r *http.Request) {
  sc := &siteCrud{h}
  h.stdcrud(w, r, sc)
}

