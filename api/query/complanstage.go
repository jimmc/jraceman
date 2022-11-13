package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type complanstageQuery struct{
  domain.ComplanStageMeta
  h *handler
}

func (sc *complanstageQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) complanstage(w http.ResponseWriter, r *http.Request) {
  sq := &complanstageQuery{domain.ComplanStageMeta{}, h}
  h.stdquery(w, r, sq)
}
