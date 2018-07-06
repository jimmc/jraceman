package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type entryQuery struct{
  h *handler
}

func (sc *entryQuery) EntityTypeName() string {
  return "entry"
}

func (sc *entryQuery) NewEntity() interface{} {
  return &domain.Entry{}
}

func (h *handler) entry(w http.ResponseWriter, r *http.Request) {
  sq := &entryQuery{h}
  h.stdquery(w, r, sq)
}
