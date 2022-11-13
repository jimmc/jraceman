package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type genderCrud struct{
  domain.GenderMeta
  h *handler
}

func (sc *genderCrud) Save(entity interface{}) (string, error) {
  var gender *domain.Gender = entity.(*domain.Gender)
  return sc.h.config.DomainRepos.Gender().Save(gender)
}

func (sc *genderCrud) List(offset, limit int) ([]interface{}, error) {
  genders, err := sc.h.config.DomainRepos.Gender().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(genders))
  for i, gender := range genders {
    a[i] = gender
  }
  return a, nil
}

func (sc *genderCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Gender().FindByID(ID)
}

func (sc *genderCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Gender().DeleteByID(ID)
}

func (sc *genderCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldGender *domain.Gender = oldEntity.(*domain.Gender)
  var newGender *domain.Gender = newEntity.(*domain.Gender)
  return sc.h.config.DomainRepos.Gender().UpdateByID(ID, oldGender, newGender, diffs)
}

func (h *handler) gender(w http.ResponseWriter, r *http.Request) {
  sc := &genderCrud{domain.GenderMeta{}, h}
  h.stdcrud(w, r, sc)
}
