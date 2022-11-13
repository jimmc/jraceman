package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type registrationfeeCrud struct{
  domain.RegistrationFeeMeta
  h *handler
}

func (sc *registrationfeeCrud) Save(entity interface{}) (string, error) {
  var registrationfee *domain.RegistrationFee = entity.(*domain.RegistrationFee)
  return sc.h.config.DomainRepos.RegistrationFee().Save(registrationfee)
}

func (sc *registrationfeeCrud) List(offset, limit int) ([]interface{}, error) {
  registrationfees, err := sc.h.config.DomainRepos.RegistrationFee().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(registrationfees))
  for i, registrationfee := range registrationfees {
    a[i] = registrationfee
  }
  return a, nil
}

func (sc *registrationfeeCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.RegistrationFee().FindByID(ID)
}

func (sc *registrationfeeCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.RegistrationFee().DeleteByID(ID)
}

func (sc *registrationfeeCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldRegistrationFee *domain.RegistrationFee = oldEntity.(*domain.RegistrationFee)
  var newRegistrationFee *domain.RegistrationFee = newEntity.(*domain.RegistrationFee)
  return sc.h.config.DomainRepos.RegistrationFee().UpdateByID(ID, oldRegistrationFee, newRegistrationFee, diffs)
}

func (h *handler) registrationfee(w http.ResponseWriter, r *http.Request) {
  sc := &registrationfeeCrud{domain.RegistrationFeeMeta{}, h}
  h.stdcrud(w, r, sc)
}
