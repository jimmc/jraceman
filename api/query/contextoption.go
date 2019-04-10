package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type contextoptionQuery struct{
  h *handler
}

func (sc *contextoptionQuery) EntityTypeName() string {
  return "contextoption"
}

func (sc *contextoptionQuery) NewEntity() interface{} {
  return &domain.ContextOption{}
}

func (sc *contextoptionQuery) SummaryQuery() string {
  return "select Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) contextoption(w http.ResponseWriter, r *http.Request) {
  sq := &contextoptionQuery{h}
  h.stdquery(w, r, sq)
}
