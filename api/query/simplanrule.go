package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type simplanruleQuery struct{
  h *handler
}

func (sc *simplanruleQuery) EntityTypeName() string {
  return "simplanrule"
}

func (sc *simplanruleQuery) NewEntity() interface{} {
  return &domain.SimplanRule{}
}

func (h *handler) simplanrule(w http.ResponseWriter, r *http.Request) {
  sq := &simplanruleQuery{h}
  h.stdquery(w, r, sq)
}
