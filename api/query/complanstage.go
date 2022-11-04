package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type complanstageQuery struct{
  h *handler
}

func (sc *complanstageQuery) EntityTypeName() string {
  return "complanstage"
}

func (sc *complanstageQuery) NewEntity() interface{} {
  return &domain.ComplanStage{}
}

func (sc *complanstageQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) complanstage(w http.ResponseWriter, r *http.Request) {
  sq := &complanstageQuery{h}
  h.stdquery(w, r, sq)
}
