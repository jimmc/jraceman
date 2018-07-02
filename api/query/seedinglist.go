package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
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

func (h *handler) seedinglist(w http.ResponseWriter, r *http.Request) {
  sq := &seedinglistQuery{h}
  h.stdquery(w, r, sq)
}
