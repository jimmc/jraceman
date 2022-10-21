package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type seedinglistQuery struct{
  h *handler
}

func (sc *seedinglistQuery) EntityTypeName() string {
  return "seedinglist"
}

func (sc *seedinglistQuery) NewEntity() interface{} {
  return &domain.SeedingList{}
}

func (sc *seedinglistQuery) SummaryQuery() string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) seedinglist(w http.ResponseWriter, r *http.Request) {
  sq := &seedinglistQuery{h}
  h.stdquery(w, r, sq)
}
