package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type registrationCrud struct{
  domain.RegistrationMeta
  h *handler
}

func (sc *registrationCrud) Save(entity interface{}) (string, error) {
  var registration *domain.Registration = entity.(*domain.Registration)
  return sc.h.config.DomainRepos.Registration().Save(registration)
}

func (sc *registrationCrud) List(offset, limit int) ([]interface{}, error) {
  registrations, err := sc.h.config.DomainRepos.Registration().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(registrations))
  for i, registration := range registrations {
    a[i] = registration
  }
  return a, nil
}

func (sc *registrationCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Registration().FindByID(ID)
}

func (sc *registrationCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Registration().DeleteByID(ID)
}

func (sc *registrationCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldRegistration *domain.Registration = oldEntity.(*domain.Registration)
  var newRegistration *domain.Registration = newEntity.(*domain.Registration)
  return sc.h.config.DomainRepos.Registration().UpdateByID(ID, oldRegistration, newRegistration, diffs)
}

func (h *handler) registration(w http.ResponseWriter, r *http.Request) {
  sc := &registrationCrud{domain.RegistrationMeta{}, h}
  h.stdcrud(w, r, sc)
}
