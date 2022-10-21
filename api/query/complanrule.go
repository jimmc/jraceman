package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type complanruleQuery struct{
  h *handler
}

func (sc *complanruleQuery) EntityTypeName() string {
  return "complanrule"
}

func (sc *complanruleQuery) NewEntity() interface{} {
  return &domain.ComplanRule{}
}

func (sc *complanruleQuery) SummaryQuery() string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) complanrule(w http.ResponseWriter, r *http.Request) {
  sq := &complanruleQuery{h}
  h.stdquery(w, r, sq)
}
