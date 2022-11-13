package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type simplanruleQuery struct{
  domain.SimplanRuleMeta
  h *handler
}

func (sc *simplanruleQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) simplanrule(w http.ResponseWriter, r *http.Request) {
  sq := &simplanruleQuery{domain.SimplanRuleMeta{}, h}
  h.stdquery(w, r, sq)
}
