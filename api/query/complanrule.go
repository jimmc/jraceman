package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type complanruleQuery struct{
  domain.ComplanRuleMeta
  h *handler
}

func (sc *complanruleQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) complanrule(w http.ResponseWriter, r *http.Request) {
  sq := &complanruleQuery{domain.ComplanRuleMeta{}, h}
  h.stdquery(w, r, sq)
}
