package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type personQuery struct{
  domain.PersonMeta
  h *handler
}

func (sc *personQuery) SummaryQuery(format string) string {
  return "select person.ID as ID, person.LastName || ', ' || person.FirstName || " +
          "' (' || team.ShortName || ') [' || person.ID || ']' as summary " +
          "from person join team on person.teamid = team.id"
}

func (h *handler) person(w http.ResponseWriter, r *http.Request) {
  sq := &personQuery{domain.PersonMeta{}, h}
  h.stdquery(w, r, sq)
}
