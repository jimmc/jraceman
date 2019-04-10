package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type raceQuery struct{
  h *handler
}

func (sc *raceQuery) EntityTypeName() string {
  return "race"
}

func (sc *raceQuery) NewEntity() interface{} {
  return &domain.Race{}
}

func (sc *raceQuery) SummaryQuery() string {
  return "select '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) race(w http.ResponseWriter, r *http.Request) {
  sq := &raceQuery{h}
  h.stdquery(w, r, sq)
}
