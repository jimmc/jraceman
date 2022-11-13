package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type siteCrud struct{
  domain.SiteMeta
  h *handler
}

func (sc *siteCrud) Save(entity interface{}) (string, error) {
  var site *domain.Site = entity.(*domain.Site)
  return sc.h.config.DomainRepos.Site().Save(site)
}

func (sc *siteCrud) List(offset, limit int) ([]interface{}, error) {
  sites, err := sc.h.config.DomainRepos.Site().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(sites))
  for i, site := range sites {
    a[i] = site
  }
  return a, nil
}

func (sc *siteCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Site().FindByID(ID)
}

func (sc *siteCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Site().DeleteByID(ID)
}

func (sc *siteCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldSite *domain.Site = oldEntity.(*domain.Site)
  var newSite *domain.Site = newEntity.(*domain.Site)
  return sc.h.config.DomainRepos.Site().UpdateByID(ID, oldSite, newSite, diffs)
}

func (h *handler) site(w http.ResponseWriter, r *http.Request) {
  sc := &siteCrud{domain.SiteMeta{}, h}
  h.stdcrud(w, r, sc)
}
