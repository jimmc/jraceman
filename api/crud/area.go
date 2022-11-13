package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type areaCrud struct{
  domain.AreaMeta
  h *handler
}

func (sc *areaCrud) Save(entity interface{}) (string, error) {
  var area *domain.Area = entity.(*domain.Area)
  return sc.h.config.DomainRepos.Area().Save(area)
}

func (sc *areaCrud) List(offset, limit int) ([]interface{}, error) {
  areas, err := sc.h.config.DomainRepos.Area().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(areas))
  for i, area := range areas {
    a[i] = area
  }
  return a, nil
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
  sc := &areaCrud{domain.AreaMeta{}, h}
  h.stdcrud(w, r, sc)
}
