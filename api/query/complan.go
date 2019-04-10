package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type complanQuery struct{
  h *handler
}

func (sc *complanQuery) EntityTypeName() string {
  return "complan"
}

func (sc *complanQuery) NewEntity() interface{} {
  return &domain.Complan{}
}

func (sc *complanQuery) SummaryQuery() string {
  return "select '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) complan(w http.ResponseWriter, r *http.Request) {
  sq := &complanQuery{h}
  h.stdquery(w, r, sq)
}
