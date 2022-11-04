package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
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

func (sc *complanQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) complan(w http.ResponseWriter, r *http.Request) {
  sq := &complanQuery{h}
  h.stdquery(w, r, sq)
}
