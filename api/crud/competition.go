package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type competitionCrud struct{
  h *handler
}

func (sc *competitionCrud) EntityTypeName() string {
  return "competition"
}

func (sc *competitionCrud) NewEntity() interface{} {
  return &domain.Competition{}
}

func (sc *competitionCrud) Save(entity interface{}) (string, error) {
  var competition *domain.Competition = entity.(*domain.Competition)
  return sc.h.config.DomainRepos.Competition().Save(competition)
}

func (sc *competitionCrud) List(offset, limit int) ([]interface{}, error) {
  competitions, err := sc.h.config.DomainRepos.Competition().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(competitions))
  for i, competition := range competitions {
    a[i] = competition
  }
  return a, nil
}

func (sc *competitionCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Competition().FindByID(ID)
}

func (sc *competitionCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Competition().DeleteByID(ID)
}

func (sc *competitionCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldCompetition *domain.Competition = oldEntity.(*domain.Competition)
  var newCompetition *domain.Competition = newEntity.(*domain.Competition)
  return sc.h.config.DomainRepos.Competition().UpdateByID(ID, oldCompetition, newCompetition, diffs)
}

func (h *handler) competition(w http.ResponseWriter, r *http.Request) {
  sc := &competitionCrud{h}
  h.stdcrud(w, r, sc)
}
