package crud

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type teamCrud struct{
  h *handler
}

func (sc *teamCrud) EntityTypeName() string {
  return "team"
}

func (sc *teamCrud) NewEntity() interface{} {
  return &domain.Team{}
}

func (sc *teamCrud) Save(entity interface{}) (string, error) {
  var team *domain.Team = entity.(*domain.Team)
  return sc.h.config.DomainRepos.Team().Save(team)
}

func (sc *teamCrud) List(offset, limit int) ([]interface{}, error) {
  teams, err := sc.h.config.DomainRepos.Team().List(offset, limit)
  if err != nil {
    return nil, err
  }
  a := make([]interface{}, len(teams))
  for i, team := range teams {
    a[i] = team
  }
  return a, nil
}

func (sc *teamCrud) FindByID(ID string) (interface{}, error) {
  return sc.h.config.DomainRepos.Team().FindByID(ID)
}

func (sc *teamCrud) DeleteByID(ID string) error {
  return sc.h.config.DomainRepos.Team().DeleteByID(ID)
}

func (sc *teamCrud) UpdateByID(ID string, oldEntity, newEntity interface{}, diffs domain.Diffs) error {
  var oldTeam *domain.Team = oldEntity.(*domain.Team)
  var newTeam *domain.Team = newEntity.(*domain.Team)
  return sc.h.config.DomainRepos.Team().UpdateByID(ID, oldTeam, newTeam, diffs)
}

func (h *handler) team(w http.ResponseWriter, r *http.Request) {
  sc := &teamCrud{h}
  h.stdcrud(w, r, sc)
}
