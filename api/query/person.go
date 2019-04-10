package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type personQuery struct{
  h *handler
}

func (sc *personQuery) EntityTypeName() string {
  return "person"
}

func (sc *personQuery) NewEntity() interface{} {
  return &domain.Person{}
}

func (sc *personQuery) SummaryQuery() string {
  return "select person.LastName || ', ' || person.FirstName || " +
          "' (' || team.ShortName || ') [' || person.ID || ']' as summary " +
          "from person join team on person.teamid = team.id"
}

func (h *handler) person(w http.ResponseWriter, r *http.Request) {
  sq := &personQuery{h}
  h.stdquery(w, r, sq)
}
