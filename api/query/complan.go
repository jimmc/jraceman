package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type complanQuery struct{
  domain.ComplanMeta
  h *handler
}

func (sc *complanQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) complan(w http.ResponseWriter, r *http.Request) {
  sq := &complanQuery{domain.ComplanMeta{}, h}
  h.stdquery(w, r, sq)
}
