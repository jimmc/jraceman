package crud

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type challengeCrud struct{
  domain.ChallengeMeta
  h *handler
}

func (sc *challengeCrud) Save(entity interface{}) (string, error) {
  var challenge *domain.Challenge = entity.(*domain.Challenge)
  return sc.h.config.DomainRepos.Challenge().Save(challenge)
}

func (sc *challengeCrud) List(offset, limit int) ([]interface{}, error) {
  challenges, err := sc.h.config.DomainRepos.Challenge().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(challenges))
  for i, challenge := range challenges {
    a[i] = challenge
  }
  return a, nil
}

func (sc *challengeCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Challenge().FindByID(ID)
}

func (sc *challengeCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Challenge().DeleteByID(ID)
}

func (sc *challengeCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldChallenge *domain.Challenge = oldEntity.(*domain.Challenge)
  var newChallenge *domain.Challenge = newEntity.(*domain.Challenge)
  return sc.h.config.DomainRepos.Challenge().UpdateByID(ID, oldChallenge, newChallenge, diffs)
}

func (h *handler) challenge(w http.ResponseWriter, r *http.Request) {
  sc := &challengeCrud{domain.ChallengeMeta{}, h}
  h.stdcrud(w, r, sc)
}
