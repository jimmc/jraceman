package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type areaCrud struct{
  h *handler
}

func (sc *areaCrud) EntityTypeName() string {
  return "area"
}

func (sc *areaCrud) NewEntity() interface{} {
  return &domain.Area{}
}

func (sc *areaCrud) Save(entity interface{}) error {
  var area *domain.Area = entity.(*domain.Area)
  return sc.h.config.DomainRepos.Area().Save(area)
}

func (sc *areaCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Area().FindByID(ID)
}

func (sc *areaCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Area().DeleteByID(ID)
}

func (sc *areaCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldArea *domain.Area = oldEntity.(*domain.Area)
  var newArea *domain.Area = newEntity.(*domain.Area)
  return sc.h.config.DomainRepos.Area().UpdateByID(ID, oldArea, newArea, diffs)
}

func (h *handler) area(w http.ResponseWriter, r *http.Request) {
  sc := &areaCrud{h}
  h.stdcrud(w, r, sc)
}
