package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type exceptionQuery struct{
  h *handler
}

func (sc *exceptionQuery) EntityTypeName() string {
  return "exception"
}

func (sc *exceptionQuery) NewEntity() interface{} {
  return &domain.Exception{}
}

func (sc *exceptionQuery) SummaryQuery() string {
  return "select Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) exception(w http.ResponseWriter, r *http.Request) {
  sq := &exceptionQuery{h}
  h.stdquery(w, r, sq)
}
