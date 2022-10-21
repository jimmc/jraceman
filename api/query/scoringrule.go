package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type scoringruleQuery struct{
  h *handler
}

func (sc *scoringruleQuery) EntityTypeName() string {
  return "scoringrule"
}

func (sc *scoringruleQuery) NewEntity() interface{} {
  return &domain.ScoringRule{}
}

func (sc *scoringruleQuery) SummaryQuery() string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) scoringrule(w http.ResponseWriter, r *http.Request) {
  sq := &scoringruleQuery{h}
  h.stdquery(w, r, sq)
}
