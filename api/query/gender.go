package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type genderQuery struct{
  h *handler
}

func (sc *genderQuery) EntityTypeName() string {
  return "gender"
}

func (sc *genderQuery) NewEntity() interface{} {
  return &domain.Gender{}
}

func (sc *genderQuery) SummaryQuery() string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) gender(w http.ResponseWriter, r *http.Request) {
  sq := &genderQuery{h}
  h.stdquery(w, r, sq)
}
