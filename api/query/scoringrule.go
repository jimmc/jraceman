package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type scoringruleQuery struct{
  domain.ScoringRuleMeta
  h *handler
}

func (sc *scoringruleQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) scoringrule(w http.ResponseWriter, r *http.Request) {
  sq := &scoringruleQuery{domain.ScoringRuleMeta{}, h}
  h.stdquery(w, r, sq)
}
